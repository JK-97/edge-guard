package register

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"jxcore/core/device"
	"jxcore/core/heartbeat"
	"jxcore/internal/network/dns"
	"jxcore/internal/network/vpn"
	"jxcore/oplog"
	"jxcore/oplog/logs"
	"jxcore/oplog/types"
	"jxcore/version"

	"jxcore/lowapi/logger"
)

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

		masterip := retryfindMaster(ctx)
		if !once {
			once = true
			onFirstConnect()
		}
		onMasterIPChanged(masterip)
		err := heartbeat.AliveReport(ctx, masterip, 5)
		if err != nil {
			logger.Error(err)
			oplog.Insert(logs.NewOplog(types.NETWORKE, fmt.Sprintf("heartbeat failed ->%s", masterip)))
		}
		oplog.Insert(logs.NewOplog(types.NETWORKE, fmt.Sprintf("heartbeat success ->%s", masterip)))
	}
}

// retryfindMaster 从 DHCP 服务器 获取 Master 节点的 IP，直到获取成功
func retryfindMaster(ctx context.Context) string {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ""
		case <-ticker.C:
			logger.Info("Try to connect a new master")
			masterip, err := findMasterFromDHCPServer(ctx)
			if err != nil {
				logger.Error("Failed to connect master: ", err)
			} else {
				return masterip
			}
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
