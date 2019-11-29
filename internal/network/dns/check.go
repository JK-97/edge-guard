package dns

import (
	"io/ioutil"
	"jxcore/core/device"
	"jxcore/lowapi/utils"
	"strings"
)

// CheckDnsmasqConf 检查 dnsmasq 的 hosts 文件
func CheckDnsmasqConf() bool {
	flag := 0
	currentdeive, err := device.GetDevice()
	utils.CheckErr(err)
	rawData, err := ioutil.ReadFile(DnsmasqHostFile)
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
