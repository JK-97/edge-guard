package register

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"jxcore/config/yaml"
	"jxcore/core/device"
	"jxcore/core/heartbeat"
	"jxcore/internal/network/dns"
	"jxcore/internal/network/iface"
	"jxcore/internal/network/vpn"
	"jxcore/version"

	"jxcore/lowapi/logger"
)

const (
	// FallBackAuthHost 默认集群地址
	FallBackAuthHost = "http://auth.iotedge.jiangxingai.com:1054"

	wireguardRegisterPath = "/api/v1/wg/register"
	openvpnRegisterPath   = "/api/v1/openvpn/register"

	prefix = 512
	suffix = 128

	registerTimeout = time.Second * 10

	consulConfigPath = "/data/edgex/consul/config/consul_conf.json"
)

var enc = base64.NewEncoding("ABCDEFGHIJKLMNOabcdefghijklmnopqrstuvwxyzPQRSTUVWXYZ0123456789-_").WithPadding(base64.NoPadding)

type reqRegister struct {
	WorkerID string `json:"wid"`
	Nonce    int64  `json:"nonce"`
	Version  string `json:"version"`
	Key      string `json:"key"`
}

type registerInfo struct {
	masterip  string
	vpnConfig []byte
}

// 维持到 IoTEdge Master的连接，第一次连接成功后执行onFirstConnect
func MaintainMasterConnection(ctx context.Context, onFirstConnect func()) error {
	once := false
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		heartbeatIP := ""
		masterPublicIP, masterVpnIP := retryfindMaster(ctx)
		if !once {
			once = true
			onFirstConnect()
		}

		if yaml.Config.HeartBeatThroughVpn {
			heartbeatIP = masterVpnIP
		} else {
			heartbeatIP = masterPublicIP
		}
		// add Public network route
		route, err := iface.GetGWRoute(iface.GetCurrentIFcae())
		if err != nil {
			logger.Error("Failed to get gwroute")
			continue
		}
		iface.SetHighPriority(route)
		err = iface.ReplcaeRouteMask32(route, heartbeatIP)
		if err != nil {
			logger.Error("Failed to add materIp route")
			continue
		}

		onMasterIPChanged(heartbeatIP)
		err = heartbeat.AliveReport(ctx, heartbeatIP, 5)
		if err != nil {
			logger.Error(err)
		}
	}
}

// retryfindMaster 从 DHCP 服务器 获取 Master 节点的 IP,vpnIp，直到获取成功
func retryfindMaster(ctx context.Context) (string, string) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ""
		case <-ticker.C:
			logger.Info("Try to connect a new master")
			// get vpnIP
			masterVpnIp, err := findMasterFromDHCPServer(ctx)
			if err != nil {
				logger.Error(err)
				continue
			}
			//get Public network IP
			masterIP, err := vpn.ParseMasterIPFromVpnConfig()
			if err != nil {
				logger.Error(err)
				continue
			}
			return masterIP, masterVpnIp

		}
	}
}

// findMasterFromDHCPServer 尝试从 DHCP 服务器 获取 Master 节点的 IP
func findMasterFromDHCPServer(ctx context.Context) (masterip string, err error) {
	var dev *device.Device
	dev, err = device.GetDevice()
	if err != nil {
		return
	}
	info, err := register(dev)
	if err != nil {
		return
	}
	masterip = info.masterip
	err = vpn.UpdateConfig(ctx, info.vpnConfig)
	if err != nil {
		return
	}
	err = dns.SaveMasterIPToHostFile(masterip)
	return
}

// 远程调用dhcpserver注册节点，获取masterip,vpnkey
func register(dev *device.Device) (*registerInfo, error) {
	reqinfo := reqRegister{
		WorkerID: dev.WorkerID,
		Nonce:    time.Now().Unix(),
		Key:      dev.Key,
		Version:  version.Version,
	}
	reqData, err := encodeReqData(reqinfo)
	if err != nil {
		return nil, err
	}

	// register vpn对应url
	var url string
	switch dev.Vpn {
	case device.VPNModeWG:
		url = dev.DhcpServer + wireguardRegisterPath
	case device.VPNModeOPENVPN:
		url = dev.DhcpServer + openvpnRegisterPath
	default:
		return nil, fmt.Errorf("VPN mode not supported: %v", dev.Vpn)
	}

	// 发送请求
	client := http.Client{Timeout: registerTimeout}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqData))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("register to master failed, status code: %s", resp.Status)
	}

	masterip := resp.Header.Get("X-Master-IP")

	// 获得加密wgkey zip 文件
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	vpnConfig, err := decodeVpnConfig(buff)
	if err != nil {
		return nil, err
	}

	info := &registerInfo{
		masterip:  masterip,
		vpnConfig: vpnConfig,
	}
	return info, nil
}

// 构造register 请求体
func encodeReqData(reqinfo reqRegister) ([]byte, error) {
	reqdata, err := json.Marshal(reqinfo)
	if err != nil {
		return nil, err
	}
	// req base64加密
	n := enc.EncodedLen(len(reqdata))
	dst := make([]byte, n)
	enc.Encode(dst, reqdata)
	return dst, nil
}

// decode vpn config from register response
func decodeVpnConfig(buff []byte) ([]byte, error) {
	r := chaosReader{
		Bytes:  buff,
		Offset: prefix,
	}
	content := make([]byte, len(buff)-prefix-suffix)
	_, err := r.read(content)
	if err != nil {
		return nil, err
	}
	return content, nil
}
