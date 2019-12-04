package dns

import (
	"fmt"
	"io/ioutil"
	"jxcore/internal/network"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"jxcore/lowapi/logger"
)

func AddMasterDns(domain string) error {
	ipRecords, err := net.LookupIP(domain)
	if err != nil {
		return err
	}

	// Shuffle 打乱 DNS 记录
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ipRecords), func(i, j int) {
		ipRecords[i], ipRecords[j] = ipRecords[j], ipRecords[i]
	})

	f, err := os.OpenFile("/etc/dnsmasq.d/dnsmasq.conf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, ip := range ipRecords {
		_, err := f.WriteString("server=/.iotedge/" + ip.String() + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// 重设 /etc/dnsmasq.hosts
func ResetHostFile() error {
	content := fmt.Sprintf("%s %s\n", network.DockerHostIP, LocalHostName)
	content += fmt.Sprintf("%s %s\n", network.DockerHostIP, IotedgeHostName)
	err := ioutil.WriteFile(DnsmasqHostFile, []byte(content), 0644)
	if err != nil {
		return err
	}
	ReloadDnsmasq()
	return nil
}

// 添加master ip 到 /etc/dnsmasq.hosts
func SaveMasterIPToHostFile(masterip string) error {
	buf, err := ioutil.ReadFile(DnsmasqHostFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	lines := strings.Split(string(buf), "\n")
	masterLine := masterip + " " + MasterHostName
	var changedLine bool
	for i, line := range lines {
		if strings.Contains(line, MasterHostName) {
			lines[i] = masterLine
			changedLine = true
			break
		}
	}
	if !changedLine {
		lines = append(lines, masterLine)
	}

	err = ioutil.WriteFile(DnsmasqHostFile, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}
	ReloadDnsmasq()
	logger.Info("Saved master ip to dnsmasq")
	return nil
}

// LoadMasterFromHostFile 从 /etc/dnsmasq.hosts 文件获取 Master 节点的 IP
func LoadMasterFromHostFile() (string, error) {
	text, err := ioutil.ReadFile(DnsmasqHostFile)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(text), "\n") {
		if strings.Contains(line, MasterHostName) {
			arr := strings.Split(line, MasterHostName)
			return strings.TrimSpace(arr[0]), err
		}
	}
	return "", fmt.Errorf("Master IP not found from %s", DnsmasqHostFile)
}
