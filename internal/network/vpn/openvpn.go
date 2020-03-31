package vpn

import (
	"context"
	"errors"
	"github.com/JK-97/edge-guard/internal/network"
	"github.com/JK-97/edge-guard/lowapi/logger"
	"github.com/JK-97/edge-guard/lowapi/system"
	"github.com/JK-97/edge-guard/lowapi/utils"
	"io/ioutil"
	"os"
	"strings"
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

func GetOpenvpnConfig() string {
	return openvpnConfigPath
}

func ParseOpenvpnConfig() (string, error) {
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
