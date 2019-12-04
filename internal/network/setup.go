package network

import (
	"os/exec"

	log "jxcore/lowapi/logger"
)

func DisableNetworkManager() {
	_ = exec.Command("systemctl", "disable", "NetworkManager").Run()
	_ = exec.Command("systemctl", "stop", "NetworkManager").Run()
	log.Info("Disabled NetworkManager")
}

// TODO 优化jxcore启动时间：disable 耗时0m0.674s, stop 耗时0m0.055s
func DisableSystemdResolved() {
	_ = exec.Command("systemctl", "disable", "systemd-resolved").Run()
	_ = exec.Command("systemctl", "stop", "systemd-resolved").Run()
	log.Info("Disabled systemd-resolved")
}
