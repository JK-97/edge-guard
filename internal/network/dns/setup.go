package dns

import (
	"io/ioutil"
	"jxcore/core/device"
	"jxcore/internal/network"
	"jxcore/internal/network/dns/dnsfile"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"jxcore/lowapi/logger"
	log "jxcore/lowapi/logger"
	"jxcore/lowapi/utils"
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
		if words[0] != "#" {
			HostRecordMap[words[len(words)-1]] = words[0]
		}

	}

	HostRecordMap[dev.WorkerID] = "127.0.0.1"
	HostRecordMap["edgex-mongo"] = "127.0.0.1"

	f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for k, v := range HostRecordMap {
		_, _ = f.WriteString(v + " " + k + "\n")
	}
	logger.Info("Appended worker id to /etc/hosts")
}

func AddMasterDns() error {
	if utils.FileExists(dnsmasqUpstreamPath) {
		return nil
	}

	device, err := device.GetDevice()
	if err != nil {
		return err
	}

	domain := network.GetHost(device.DhcpServer)
	ipRecords, err := net.LookupIP(domain)
	if err != nil {
		return err
	}

	// Shuffle 打乱 DNS 记录
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ipRecords), func(i, j int) {
		ipRecords[i], ipRecords[j] = ipRecords[j], ipRecords[i]
	})

	f, err := os.OpenFile(dnsmasqUpstreamPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
