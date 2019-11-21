package dns

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jxcore/core/device"
	"jxcore/internal/network"
	"jxcore/lowapi/utils"
	"jxcore/template"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
)

func LookUpDns(domain string) {
	ipRecords, _ := net.LookupIP(domain)
	Shuffle(ipRecords)
	f, err := os.OpenFile("/etc/dnsmasq.d/dnsmasq.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Info(err)
		return
	}
	defer f.Close()
	for _, ip := range ipRecords {
		f.WriteString("server=/.iotedge/" + ip.String() + "\n")
	}
}

// Shuffle 打乱 DNS 记录
func Shuffle(slice []net.IP) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func RestartDnsmasq() {
	exec.Command("/bin/bash", "-c", "systemctl restart dnsmasq").Run()
}

func ResetHostFile(ethIp string) error {
	content := fmt.Sprintf("%s %s\n", ethIp, LocalHostName)
	content += fmt.Sprintf("%s %s\n", ethIp, IotedgeHostName)
	return ioutil.WriteFile(HostsFile, []byte(content), 0644)
}

func UpdateMasterIPToHosts(masterip string) {
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

	RestartDnsmasq()
}

// FindMasterFromHostFile 从 hosts 文件获取 Master 节点的 IP
func FindMasterFromHostFile() string {
	f, err := os.Open(HostsFile)
	if err != nil {
		utils.CheckErr(err)
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

// OnVPNConnetced VPN 连接成功后执行
func OnVPNConnetced() {
}

// updateConsulConfig 更新 Consul 配置
func updateConsulConfig(currentdevice *device.Device) {
	config := consulConfig{
		Server:           true,
		ClientAddr:       "0.0.0.0",
		AdvertiseAddrWan: network.GetClusterIP(),
		BootstrapExpect:  1,
		Datacenter:       "worker-" + currentdevice.WorkerID,
		NodeName:         "worker-" + currentdevice.WorkerID,
		RetryJoinWan:     []string{MasterHostName},
		UI:               true,
	}
	if buf, err := json.Marshal(config); err == nil {
		ioutil.WriteFile(consulConfigPath, buf, 0666)
	}
}

// updateTelegrafConfig 更新 Telegraf 和 InfluxDB ,cadvisor配置
func updateTelegrafConfig(currentdevice *device.Device, masterip string) {
	template.Telegrafcfg(masterip, currentdevice.WorkerID)
	template.Cadvisorcfg(masterip, currentdevice.WorkerID)
	var VpnIP string
	//确保4g 或 以太有一个起来的情况下
	if _, erreth0 := network.GetMyIP("eth0"); erreth0 == nil {
		VpnIP = network.GetClusterIP()
	} else if _, errusb0 := network.GetMyIP("usb0"); errusb0 == nil {
		VpnIP = network.GetClusterIP()
	}
	if VpnIP != "" {
		template.Statsitecfg(masterip, VpnIP)
	}
}

// OnMasterIPChanged master IP 变化后执行
func OnMasterIPChanged(masterip string) {
	currentdevice, err := device.GetDevice()
	utils.CheckErr(err)

	if utils.Exists(consulConfigPath) {
		updateConsulConfig(currentdevice)
	}

	updateTelegrafConfig(currentdevice, masterip)
}

// AppendHostnameHosts 将更新后的 hostsname 写入 hosts 文件
func AppendHostnameHosts(workerid string) {
	hostnameRecord := "127.0.0.1 " + workerid
	content, err := ioutil.ReadFile(HostsFile)
	if err != nil {
		log.Error(err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(content))

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), hostnameRecord) {
			return
		}
	}

	var hostnameresolv = "\n" + hostnameRecord + "\n127.0.0.1 " + workerid
	f, err := os.OpenFile(HostsFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Error(err)
		return
	}
	f.WriteString(hostnameresolv)
	f.Close()
}

func ParseIpInTxt(url string) (string, string) {
	txtRecords, err := net.LookupTXT(url)
	if err != nil {
		log.Error(err)
		log.Info("Possible DNS configuration error")
	}
	if len(txtRecords) == 0 {
		return "", ""
	}
	res := strings.Split(txtRecords[0], ":")
	return res[0], res[1]
}

// CheckDnsmasqConf 检查 dnsmasq 的 hosts 文件
func CheckDnsmasqConf() bool {
	flag := 0
	currentdeive, err := device.GetDevice()
	utils.CheckErr(err)
	rawData, err := ioutil.ReadFile(HostsFile)
	utils.CheckErr(err)
	lines := strings.Split(string(rawData), "\n")
	for _, line := range lines {
		if strings.Contains(line, MasterHostName) {
			flag++
		} else if strings.Contains(line, IotedgeHostName) {
			flag++
		} else if strings.Contains(line, LocalHostName) {
			flag++
		} else if strings.Contains(line, "worker-"+currentdeive.WorkerID) {
			flag++
		}
	}
	return flag >= 3

}
