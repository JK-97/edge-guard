package vpn

import (
	"context"
	"jxcore/internal/network"
	"jxcore/lowapi/system"
	"jxcore/lowapi/utils"
	"os"
	"time"
)

const (
	OpenVPNInterface = "tun0"
	openvpnConfigDir = "/etc/openvpn/"

	// openvpnSuccessMessage = "Initialization Sequence Completed"
	openvpnConfigPath  = "/etc/openvpn/client.ovpn"
	openvpnConfigName  = "iotedge"
	openvpnSoftLink    = "/etc/openvpn/" + openvpnConfigName + ".conf"
	openvpnServiceName = "openvpn@" + openvpnConfigName
)

type openvpn struct{}

// 尝试设置vpn配置
func (v *openvpn) setup() error {
	if !utils.FileExists(openvpnSoftLink) {
		err := os.Symlink(openvpnConfigPath, openvpnSoftLink)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *openvpn) retryStart(ctx context.Context) error {
	utils.RunAndLogError(v.setup)
	err := utils.RunUntilSuccess(ctx, time.Second, func() error { return system.RunCommand("systemctl enable " + openvpnServiceName) })
	if err != nil {
		return err
	}
	return utils.RunUntilSuccess(ctx, time.Second, func() error { return system.RunCommand("systemctl start " + openvpnServiceName) })
}

func (v *openvpn) stop(ctx context.Context) error {
	err := system.RunCommand("systemctl disable " + openvpnServiceName)
	if err != nil {
		return err
	}
	return system.RunCommand("systemctl stop " + openvpnServiceName)
}

func (v *openvpn) updateConfig(vpnConfig []byte) error {
	return utils.Unzip(vpnConfig, openvpnConfigDir)
}

func (v *openvpn) getIp(ctx context.Context) (string, error) {
	var ip string
	err := utils.RunUntilSuccess(ctx, time.Second*3, func() error {
		var err error
		ip, err = network.GetMyIP(OpenVPNInterface)
		return err
	})
	return ip, err
}
