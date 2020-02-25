package vpn

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"jxcore/core/device"
	"jxcore/lowapi/logger"
	"strings"
	"time"
)

const (
	vpnStartTimeout = 20 * time.Second
)

var (
	lastVpnIp   string
	vpnInstance vpn = nil

	ErrMasterNotConnected = fmt.Errorf("master not connected")
)

type vpn interface {
	// start 尝试打开 vpn 直到成功或超时
	retryStart(context.Context) error
	// stop 停止vpn
	stop(context.Context) error
	// 更新配置文件
	updateConfig([]byte) error
	// 获取节点的vpn ip
	getIp(context.Context) (string, error)
}

func getVpnInstance() (vpn, error) {
	if vpnInstance != nil {
		return vpnInstance, nil
	}
	dev, err := device.GetDevice()
	if err != nil {
		return nil, err
	}
	mode := dev.Vpn
	switch mode {
	case device.VPNModeOPENVPN:
		vpnInstance = &openvpn{}
	case device.VPNModeWG:
		vpnInstance = &wireguard{}
	default:
		return nil, fmt.Errorf("VPN mode not supported: %v", mode)
	}
	return vpnInstance, nil
}

// 更新vpn配置
func UpdateConfig(ctx context.Context, vpnConfig []byte) error {
	logger.Info("Updating VPN")
	vpn, err := getVpnInstance()
	if err != nil {
		return err
	}

	if err := vpn.updateConfig(vpnConfig); err != nil {
		return err
	}

	return Restart(ctx)
}

func Restart(ctx context.Context) error {
	vpn, err := getVpnInstance()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, vpnStartTimeout)
	defer cancel()

	if err := vpn.stop(ctx); err != nil {
		return err
	}

	if err := vpn.retryStart(ctx); err != nil {
		return err
	}

	ip, err := vpn.getIp(ctx)
	if err != nil {
		return fmt.Errorf("Failed to update cluster ip: %w", err)
	}
	lastVpnIp = ip

	return err
}

// GetClusterIP 获取集群内网 VPN IP
func GetClusterIP() string {
	return lastVpnIp
}

func ParseMasterIPFromVpnConfig() (string, error) {
	dev, err := device.GetDevice()
	if err != nil {
		return "", fmt.Errorf("Can not get device config")

	}
	switch dev.Vpn {
	case device.VPNModeOPENVPN:
		return parseOpenvpnConfig()
	case device.VPNModeWG:
		return "", fmt.Errorf("cant not support wg")
	}
	return "", fmt.Errorf("cant not support %v", dev.Vpn)
}
func parseOpenvpnConfig() (string, error) {
	data, err := ioutil.ReadFile(openvpnConfigPath)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "remote") {
			res := strings.Split(line, " ")
			logger.Info(res)
			if len(res) > 2 {
				return res[1], nil
			}
		}
	}
	return "", errors.New("parse config failed")
}
