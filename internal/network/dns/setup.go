package dns

import (
	"io/ioutil"
	"jxcore/core/device"
	"jxcore/internal/network/dns/dnsfile"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"jxcore/lowapi/logger"
	log "jxcore/lowapi/logger"
)

// 添加dhcp hook，使得dhclient的resolv.conf 结果重定向到 /edge/resolv.d/dhclient.$interface
func ApplyDHCPResolveUpdateHooks() {
	dir := filepath.Dir(dhclientResolvHookPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0644)
		if err != nil {
			logger.Fatal(err)
		}
	}

	err := ioutil.WriteFile(dhclientResolvHookPath, []byte(dnsfile.DhclientResolvRedirectHook), 0644)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Redirect dhclient resolv to /edge/resolv.d/dhclient.$interface")
}

func RemoveDHCPEnterHooks() {
	bakDir := dhcpEnterHooksDir + ".bak"
	if _, err := os.Stat(bakDir); os.IsNotExist(err) {
		_ = os.Rename(dhcpEnterHooksDir, bakDir)
		logger.Info("Removed dhcp enter hooks")
	}
}

// 重设/etc/resolv.conf 指向 dnsmasq
func ResetResolv() {
	_ = os.Remove(resolvPath)
	err := ioutil.WriteFile(resolvPath, []byte(dnsfile.ResolvConf), 0777)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Reset ", resolvPath)
}

// 重置/etc/dnsmasq.conf
func ResetDNSMasqConf() {
	ensureDNSMasqIface()
	err := ioutil.WriteFile(dnsmasqConfPath, []byte(dnsfile.DNSMasqConf), 0777)
	if err != nil {
		logger.Fatal(err)
	}
	RestartDnsmasq()
	logger.Info("Reset ", dnsmasqConfPath)
}

const dnsmasqIfaceConf = `
auto lo:0
iface lo:0 inet static
address 172.18.1.1
`

func ensureDNSMasqIface() {
	data, err := ioutil.ReadFile(ifacePath)
	if err != nil {
		logger.Fatal(err)
	}
	if !strings.Contains(string(data), "lo:0") {
		err = ioutil.WriteFile(ifacePath, append(data, []byte(dnsmasqIfaceConf)...), 0777)
		if err != nil {
			logger.Fatal(err)
		}
	}
	_ = exec.Command("ifup", "lo:0").Run()
	logger.Info("Ensured dnsmasq virtual network interface up")
}

// AppendHostnameHosts 将 127.0.0.1 workerid / mongo 加入 /etc/hosts
func AppendHostnameHosts() {
	dev, err := device.GetDevice()
	if err != nil {
		logger.Fatal(err)
	}

	content, err := ioutil.ReadFile(HostFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	data := strings.TrimSpace(string(content))

	HostRecordMap := map[string]string{}
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		trimedLine := strings.TrimSpace(line)
		words := strings.Split(trimedLine, " ")
		HostRecordMap[words[len(words)-1]] = words[0]
	}

	HostRecordMap[dev.WorkerID] = "127.0.0.1"
	HostRecordMap["mongo"] = "127.0.0.1"

	f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for k, v := range HostRecordMap {
		_, _ = f.WriteString(v + " " + k)
	}
	logger.Info("Appended worker id to /etc/hosts")
}
