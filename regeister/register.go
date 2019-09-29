package regeister

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"jxcore/log"
	"jxcore/template"
	"jxcore/utils"
	"jxcore/version"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var prefix = 512
var suffix = 128
var enc = base64.NewEncoding("ABCDEFGHIJKLMNOabcdefghijklmnopqrstuvwxyzPQRSTUVWXYZ0123456789-_").WithPadding(base64.NoPadding)

type registerRequest struct {
	WorkerID string `json:"wid"`
	Nonce    int64  `json:"nonce"`
	Version  string `json:"version"`
	Key      string `json:"key"`
}

func register(workerid, key string) {
	var err error
	var mymasterip string
	for {
		// register api 获取vpn key，链接上master
		findMasterFromDHCPServer(workerid, key)
		mymasterip, err = getmymaster(workerid, key)
		if err == nil {
			break
		}
		log.Error("Register Worker Net", err)
		// 重试
		time.Sleep(5 * time.Second)
	}
	time.Sleep(5 * time.Second)

	// 启动心跳
	AliveReport(mymasterip)
}

// getEthIP 获取以太网 IP
func getEthIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range interfaces {
		if i.Name == "eth0" {
			if addrs, err := i.Addrs(); err == nil {
				for _, addr := range addrs {
					if ipNet, ok := addr.(*net.IPNet); ok {
						return ipNet.IP.String()
					}
				}
			}

		}

	}
	return "127.0.0.1"
}

// HostsFile Hosts 文件
const HostsFile string = "/etc/dnsmasq.hosts"

// MasterHostName hosts 记录中 master 的 VPN IP
const MasterHostName string = "master.iotedge"

// updateMasterIPToHosts 更新 hosts 文件中的 master.iotedge 记录
func updateMasterIPToHosts(masterip string) {
	buf, err := ioutil.ReadFile(HostsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建文件
			f, _ := os.Create(HostsFile)
			f.Close()
			buf = make([]byte, 0)
		} else {
			log.Error(err)
			return
		}
	}
	var flag bool
	lines := strings.Split(string(buf), "\n")
	for i, line := range lines {
		if strings.Contains(line, MasterHostName) {
			lines[i] = masterip + " " + MasterHostName
			flag = true
			break
		}
	}
	var tmp []string
	if flag == false {
		tmp = append(lines, masterip+" "+MasterHostName+"\n")
	} else {
		tmp = lines
	}

	output := strings.Join(tmp, "\n")
	err = ioutil.WriteFile(HostsFile, []byte(output), 0644)
	if err != nil {
		log.Error(err)
	}
}

// findMasterFromHostFile 从 hosts 文件获取 Master 节点的 IP
func findMasterFromHostFile() string {
	f, err := os.Open(HostsFile)
	if err != nil {
		return ""
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, MasterHostName) {
			arr := strings.Split(line, MasterHostName)
			return strings.TrimSpace(arr[0])
		}
	}
	return ""
}

type consulConfig struct {
	Server           bool     `json:"server"`
	ClientAddr       string   `json:"client_addr"`
	AdvertiseAddrWan string   `json:"advertise_addr_wan"`
	BootstrapExpect  int      `json:"bootstrap_expect"`
	Datacenter       string   `json:"datacenter"`
	NodeName         string   `json:"node_name"`
	RetryJoinWan     []string `json:"retry_join_wan"`
	UI               bool     `json:"ui"`
}

const consulConfigPath = "/data/edgex/consul/config/consul_conf.json"

// onVPNConnetced VPN 连接成功后执行
func onVPNConnetced() {
	workerid := DeviceKey.WorkID
	masterip := findMasterFromHostFile()
	if masterip == "" {
		return
	}
	// 生成conf
	template.Telegrafcfg(masterip, workerid)
	var VpnIP string
	//确保4g 或 以太有一个起来的情况下
	if _, erreth0 := GetMyIP("eth0"); erreth0 == nil {
		VpnIP = GetClusterIP()
	} else if _, errusb0 := GetMyIP("usb0"); errusb0 == nil {
		VpnIP = GetClusterIP()
	}
	if VpnIP != "" {
		template.Statsitecfg(masterip, VpnIP)
	}
}

// onMasterIPChanged master IP 变化后执行
func onMasterIPChanged(masterip string) {
	workerid := DeviceKey.WorkID
	if utils.Exists(consulConfigPath) {
		config := consulConfig{
			Server:           true,
			ClientAddr:       "0.0.0.0",
			AdvertiseAddrWan: GetClusterIP(),
			BootstrapExpect:  1,
			Datacenter:       "worker-" + workerid,
			NodeName:         "worker-" + workerid,
			RetryJoinWan:     []string{MasterHostName},
			UI:               true,
		}
		if buf, err := json.Marshal(config); err == nil {
			ioutil.WriteFile(consulConfigPath, buf, 0666)
		}
	}

}

// findMasterFromDHCPServer 从 DHCP 服务器 获取 Master 节点的 IP
func findMasterFromDHCPServer(workerid string, key string) (masterip string, err error) {
	reqinfo := registerRequest{
		WorkerID: workerid,
		Nonce:    time.Now().Unix(),
		Key:      key,
		Version:  version.Version,
	}

	reqdata, err := json.Marshal(reqinfo)
	if err != nil {
		log.Error(err)
	}
	//req base64加密
	n := enc.EncodedLen(len(reqdata))
	dst := make([]byte, n)
	enc.Encode(dst, reqdata)

	//通过dhcpserver获取key
	reqbody := bytes.NewBuffer(dst)

	url := DeviceKey.DhcpServer + wireguardRegisterPath
	if DeviceKey.Vpn == VPNModeOPENVPN {
		url = DeviceKey.DhcpServer + openvpnRegisterPath
	}

	resp, err := http.Post(url, "application/json", reqbody)
	if err != nil {
		log.Error(err, "restart dnsmasq")
		exec.Command("service", "dnsmasq", "restart").Run()
		getmymaster(workerid, key)
		return
	} else if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	masterip = resp.Header.Get("X-Master-IP")
	defer resp.Body.Close()

	if masterip != "" {
		ip := findMasterFromHostFile()

		updateMasterIPToHosts(masterip)
		if masterip != ip {
			onMasterIPChanged(masterip)
		}
	}

	//获得加密wgkey zip
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	//解密
	r := ChaosReader{
		Bytes:  buff,
		Offset: prefix,
	}
	content := make([]byte, len(buff)-prefix-suffix)
	_, err = r.Read(content)

	log.Info("Updating VPN")
	// 替换vpn配置
	switch DeviceKey.Vpn {
	case VPNModeWG:
		log.Info("VPN Mode: ", DeviceKey.Vpn)
		replacesetting(bytes.NewReader(content), "/etc/wireguard")
		//vpn commponet检测配置变动 启动wireguard ,wg0
		CloseWg()
		if err := StartWg(); err == nil {
			onVPNConnetced()
		}

	case VPNModeOPENVPN:
		log.Info("VPN Mode: ", DeviceKey.Vpn)
		replacesetting(bytes.NewReader(content), "/etc/openvpn/")
		Closeopenvpn()
		if err := Startopenvpn(); err == nil {
			onVPNConnetced()
		}
	default:
		log.Error("err model")
		return
	}

	exec.Command("service", "dnsmasq", "restart").Run()
	return
}

// appendHostnameHosts 将更新后的 hostsname 写入 hosts 文件
func appendHostnameHosts(workerid string) {

	var hostnameresolv = "\n127.0.0.1 worker-" + workerid + "\n"
	f, err := os.OpenFile(HostsFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Error(err)
	}
	f.WriteString(hostnameresolv)
	f.Close()
}

func getmymaster(workerid, key string) (mymasterip string, err error) {
	masterip := findMasterFromHostFile()

	if masterip == "" {
		masterip, err = findMasterFromDHCPServer(workerid, key)
	}
	if err != nil {
		return
	}
	appendHostnameHosts(workerid)

	log.Info("Finish Update VPN")
	// _, errusb0 := GetMyIP("usb0")

	return masterip, err
}

// BuildWokerID 生成wokerid ，获取uuid
func BuildWokerID() string {
	perfilx := "J"
	if runtime.GOARCH == "amd64" {
		perfilx = perfilx + "02"
	} else {
		perfilx = perfilx + "01"
	}

	content, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		log.Error(err)
	}
	md5info := md5.Sum(content[len(content)-17:])
	md5str := fmt.Sprintf("%x", md5info)
	//获取wokerid
	workerid := perfilx + md5str[len(md5str)-7:]
	return workerid
}

func replacesetting(formfile *bytes.Reader, toetc string) {
	formfile.Seek(0, io.SeekStart)

	buff, err := ioutil.ReadAll(formfile)
	if err != nil {
		log.Error(err)
	}

	//time.Sleep(2 * time.Second)
	err = utils.Unzip(buff, toetc)
	if err != nil {
		log.Error(err)
	}
}

// ChaosReader 读取服务器混淆后的数据
type ChaosReader struct {
	Bytes  []byte
	Offset int
}

func (r *ChaosReader) Read(p []byte) (n int, err error) {
	length := len(r.Bytes)
	remain := length - r.Offset
	if remain <= 0 {
		return 0, io.EOF
	}
	length = len(p)
	if length > remain {
		err = io.EOF
	} else {
		remain = length
	}

	for n = 0; n < remain; n++ {
		b := r.Bytes[r.Offset+n]
		if b >= 0x80 {
			p[n] = b - 0x80
		} else {
			p[n] = b + 0x80
		}
	}
	r.Offset += remain
	return
}
