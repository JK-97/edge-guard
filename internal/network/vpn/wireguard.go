package vpn

import (
	"context"
	"errors"
	"io/ioutil"
	"jxcore/internal/network"
	log "jxcore/lowapi/logger"
	"jxcore/lowapi/utils"
	"net"
	"os/exec"
	"strings"
)

const (
	WireGuardInterface  = "wg0"
	wireguardConfigDir  = "/etc/wireguard/"
	wireguardConfigFile = "/etc/wireguard/wg0.conf"
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

func ParseWireGuardConfig() (string, error) {
	data, err := ioutil.ReadFile(wireguardConfigFile)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Endpoint") {
			res := strings.Split(line, " ")
			if len(res) >= 2 {
				ip, _, err := net.SplitHostPort(res[2])
				if err != nil {
					return "", err
				}
				return ip, nil
			}
		}
	}
	return "", errors.New("parse config failed")
}
