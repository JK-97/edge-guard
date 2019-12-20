package dns

import (
	"fmt"
	"io/ioutil"
	"jxcore/internal/network"
	"os"
	"strings"

	"jxcore/lowapi/logger"
)

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
