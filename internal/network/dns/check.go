package dns

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// CheckDnsmasqConf 检查 dnsmasq 的 hosts 文件
func CheckDnsmasqConf() error {
	rawData, err := ioutil.ReadFile(DnsmasqHostFile)
	if err != nil {
		return err
	}

	text := string(rawData)
	result := ""
	check := func(target string) {
		if !strings.Contains(text, target) {
			result += fmt.Sprintf("%s check failed, %s not found.\n", DnsmasqHostFile, target)
		}
	}
	check(MasterHostName)
	check(IotedgeHostName)
	check(LocalHostName)
	if result != "" {
		return fmt.Errorf(result)
	}
	return nil
}
