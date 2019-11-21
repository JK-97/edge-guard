package network

import (
	"os/exec"

	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
)

func DisableNetworkManager() {
	_ = exec.Command("systemctl", "disable", "NetworkManager").Run()
	_ = exec.Command("systemctl", "stop", "NetworkManager").Run()
	log.Info("Disabled NetworkManager")
}

func DisableSystemdResolved() {
	_ = exec.Command("systemctl", "disable", "systemd-resolved").Run()
	_ = exec.Command("systemctl", "stop", "systemd-resolved").Run()
	log.Info("Disabled systemd-resolved")
}
