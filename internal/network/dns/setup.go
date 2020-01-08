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

	"jxcore/lowapi/docker"
	"jxcore/lowapi/logger"
	"jxcore/lowapi/system"
	"jxcore/lowapi/utils"
)

// 尝试设置dns配置，忽略出错
func TrySetupDnsConfig() {
	// NetworkManager 和 systemd-resolved 会更改 /etc/resolv.conf，使dns不可控。需要停止
	utils.RunAndLogError(func() error { return system.StopDisableService("NetworkManager") })
	utils.RunAndLogError(func() error { return system.StopDisableService("systemd-resolved") })

	// 向前兼容：旧版本jxcore会锁/etc/resolv.conf
	utils.RunAndLogError(unlockResolvConf)

	// nano: /etc/dhcp/dhclient-enter-hooks.d 里的脚本需要被移除。
	// 如果脚本存在，`ifup`时不会更新 /etc/resolv.conf
	utils.RunAndLogError(removeDHCPEnterHooks)

	// 将dhcp 获取的resolv 信息重定向到 /edge/resolv.d
	utils.RunAndLogError(applyDHCPResolveUpdateHooks)
	utils.RunAndLogError(resetResolv)
	utils.RunAndLogError(resetDNSMasqConf)
	utils.RunAndLogError(appendHostnameHosts)
	utils.RunAndLogError(resetHostFile)

	utils.RunAndLogError(func() error { return AddMasterDns(false) })
	utils.RunAndLogError(docker.EnsureDockerDNSConfig)
}

// 添加dhcp hook，使得dhclient的resolv.conf 结果重定向到 /edge/resolv.d/dhclient.$interface
func applyDHCPResolveUpdateHooks() error {
	dir := filepath.Dir(dhclientResolvHookPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0644)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(dhclientResolvHookPath, []byte(dnsfile.DhclientResolvRedirectHook), 0644)
	if err != nil {
		return err
	}
	logger.Info("Redirect dhclient resolv to /edge/resolv.d/dhclient.$interface")
	return nil
}

func removeDHCPEnterHooks() error {
	bakDir := dhcpEnterHooksDir + ".bak"
	if _, err := os.Stat(bakDir); os.IsNotExist(err) {
		logger.Info("Remove dhcp enter hooks")
		return os.Rename(dhcpEnterHooksDir, bakDir)
	}
	return nil
}

// 重设/etc/resolv.conf 指向 dnsmasq
func resetResolv() error {
	_ = os.Remove(resolvPath)
	err := ioutil.WriteFile(resolvPath, []byte(dnsfile.ResolvConf), 0777)
	if err != nil {
		return err
	}
	logger.Info("Reset ", resolvPath)
	return nil
}

// 重置/etc/dnsmasq.conf
func resetDNSMasqConf() error {
	err := ensureDNSMasqIface()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dnsmasqConfPath, []byte(dnsfile.DNSMasqConf), 0777)
	if err != nil {
		return err
	}
	RestartDnsmasq()
	logger.Info("Reset ", dnsmasqConfPath)
	return nil
}

const dnsmasqIfaceConf = `
auto lo:0
iface lo:0 inet static
address 172.18.1.1
`

func ensureDNSMasqIface() error {
	data, err := ioutil.ReadFile(ifacePath)
	if err != nil {
		return err
	}
	if !strings.Contains(string(data), "lo:0") {
		err = ioutil.WriteFile(ifacePath, append(data, []byte(dnsmasqIfaceConf)...), 0777)
		if err != nil {
			return err
		}
	}
	err = exec.Command("ifup", "lo:0").Run()
	logger.Info("Ensured dnsmasq virtual network interface up")
	return err
}

// appendHostnameHosts 将 127.0.0.1 workerid / mongo 加入 /etc/hosts
func appendHostnameHosts() error {
	dev, err := device.GetDevice()
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(HostFile)
	if err != nil && !os.IsNotExist(err) {
		return err
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
		return err
	}
	defer f.Close()

	for k, v := range HostRecordMap {
		_, _ = f.WriteString(v + " " + k + "\n")
	}
	logger.Info("Appended worker id to /etc/hosts")
	return nil
}

// 向/etc/dnsmasq.d/dnsmasq.conf 添加 master dns 为本地dnsmasq 的上游服务器
func AddMasterDns(force bool) error {
	if !force && utils.FileExists(dnsmasqUpstreamPath) {
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

// backward compatible
func unlockResolvConf() error {
	return system.RunCommand("chattr -i " + resolvPath)
}
