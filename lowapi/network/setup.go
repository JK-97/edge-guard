package network

import (
	log "jxcore/go-utils/logger"
	"os"
	"os/exec"
)

const (
	dhcpEnterHooksDir = "/etc/dhcp/dhclient-enter-hooks.d"
)

// 1. NetworkManager 和 systemd-resolved 会更改 /etc/resolv.conf，使dns不可控。需要停止
// 2. /etc/dhcp/dhclient-enter-hooks.d 里的脚本需要被移除。
// 		如果脚本存在，`ifup`时不会更新 /etc/resolv.conf
func SetupNetwork() {
	disableNetworkManager()
	disableSystemdResolved()
	removeDHCPEnterHooks()
}

func disableNetworkManager() {
	log.Info("Disabling NetworkManager")
	_ = exec.Command("systemctl", "disable", "NetworkManager").Run()
	_ = exec.Command("systemctl", "stop", "NetworkManager").Run()
}

func disableSystemdResolved() {
	log.Info("Disabling systemd-resolved")
	_ = exec.Command("systemctl", "disable", "systemd-resolved").Run()
	_ = exec.Command("systemctl", "stop", "systemd-resolved").Run()
}

func removeDHCPEnterHooks() {
	err := os.Rename(dhcpEnterHooksDir, dhcpEnterHooksDir+".bak")
	if err != nil {
		log.Error(err)
	}
}
