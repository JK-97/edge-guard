package dns

import (
	"io/ioutil"
	"jxcore/internal/network/dns/dnsfile"
	"os"
	"os/exec"
	"strings"

	"gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
)

// 添加dhcp hook，使得dhclient的resolv.conf 结果重定向到 /edge/resolv.d/dhclient.$interface
func ApplyDHCPResolveUpdateHooks() {
	err := ioutil.WriteFile(dhclientResolvHookPath, []byte(dnsfile.DhclientResolvRedirectHook), 0644)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Redirect dhclient resolv to /edge/resolv.d/dhclient.$interface")
}

func RemoveDHCPEnterHooks() {
	_ = os.Rename(dhcpEnterHooksDir, dhcpEnterHooksDir+".bak")
	logger.Info("Removed dhcp enter hooks")
}

// 重设/etc/resolv.conf 指向 dnsmasq
func ResetResolv() {
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
		_ = exec.Command("ifup", "lo:0").Run()
	}
}
