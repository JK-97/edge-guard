package vpn

import (
	"context"
	"jxcore/internal/network"
	log "jxcore/lowapi/logger"
	"jxcore/lowapi/utils"
	"os/exec"
)

const (
	WireGuardInterface = "wg0"
	wireguardConfigDir = "/etc/wireguard/"
)

type wireguard struct{}

func (v *wireguard) retryStart(ctx context.Context) error {
	cmd := exec.Command("wg-quick", "up", WireGuardInterface)
	out, err := cmd.CombinedOutput()
	if err == nil {
		log.Info("wg up success")
	} else {
		log.Info("wg up failed: ", err, string(out))
	}

	return err
}

func (v *wireguard) stop(ctx context.Context) error {
	cmd := exec.Command("wg-quick", "down", WireGuardInterface)
	out, err := cmd.Output()
	if err == nil {
		log.Info("wg down success")
	} else {
		log.Info("wg down failed: ", err, string(out))
	}
	return err
}

func (v *wireguard) updateConfig(vpnConfig []byte) error {
	return utils.Unzip(vpnConfig, wireguardConfigDir)
}

func (v *wireguard) getIp(ctx context.Context) (string, error) {
	return network.GetMyIP(WireGuardInterface)
}
