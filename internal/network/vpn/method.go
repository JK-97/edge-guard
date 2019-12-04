package vpn

import (
	"bufio"
	"fmt"
	"jxcore/core/device"
	"jxcore/internal/network"
	"jxcore/lowapi/utils"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
)

// 更新vpn配置
func UpdateVPN(vpnConfig []byte) error {
	log.Info("Updating VPN")
	dev, err := device.GetDevice()
	if err != nil {
		return err
	}
	mode := dev.Vpn

	log.Info("VPN Mode: ", mode)

	var configDir string
	switch mode {
	case device.VPNModeWG:
		configDir = wireguardConfigDir
	case device.VPNModeOPENVPN:
		configDir = openvpnConfigDir
	default:
		return fmt.Errorf("VPN mode not supported: %v", mode)
	}

	if err := utils.Unzip(vpnConfig, configDir); err != nil {
		return err
	}
	if err := StopVpn(mode); err != nil {
		return err
	}
	if err := StartVpn(mode); err != nil {
		return err
	}
	return nil
}

func StartVpn(mode device.Vpn) error {
	switch mode {
	case device.VPNModeWG:
		return startWg()
	case device.VPNModeOPENVPN:
		return startopenvpn()
	default:
		return fmt.Errorf("VPN mode not supported: %v", mode)
	}
}

func StopVpn(mode device.Vpn) error {
	switch mode {
	case device.VPNModeWG:
		return stopWg()
	case device.VPNModeOPENVPN:
		return stopOpenvpn()
	default:
		return fmt.Errorf("VPN mode not supported: %v", mode)
	}
}

// startWg 打开 WireGuard VPN
func startWg() error {
	cmd := exec.Command("wg-quick", "up", WireGuardInterface)
	out, err := cmd.CombinedOutput()
	if err == nil {
		log.Info("wg up success")
	} else {
		log.Info("wg up failed: ", err, string(out))
	}

	return err
}

// stopWg 关闭 WireGuard VPN
func stopWg() error {
	cmd := exec.Command("wg-quick", "down", WireGuardInterface)
	out, err := cmd.Output()
	if err == nil {
		log.Info("wg down success")
	} else {
		log.Info("wg down failed: ", err, string(out))
	}
	return err
}

// startopenvpn 打开 OpenVPN
func startopenvpn() error {
	cmd := exec.Command("openvpn", openvpnConfigPath)
	pipe, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		log.Info("openvpn up failed :", err)
		return err
	}
	scanner := bufio.NewScanner(pipe)
	// 检测 OpenVPN 是否正常启动
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), openvpnSuccessMessage) {
				pipe.Close()
				wg.Done()
				break
			}
		}
		if scanner.Err() == nil {
			return
		}
		if perr, ok := scanner.Err().(*os.PathError); ok {
			if perr.Err != os.ErrClosed {
				err = scanner.Err()
			}
		} else {
			err = scanner.Err()
		}
	}()
	timer := time.NewTimer(openvpnWaitTimeout)
	go func() {
		select {
		case <-timer.C:
			err = ErrOpenVPNTimeout
			pipe.Close()
			return
		}
	}()
	wg.Wait()
	if err == nil {
		log.Info("openvpn up success")
	}
	return err
}

// stopOpenvpn 关闭 OpenVPN
func stopOpenvpn() error {
	c := "killall openvpn"
	err := exec.Command("/bin/sh", "-c", c).Run()
	if err == nil {
		log.Info("openvpn down success")
	} else {
		log.Info("openvpn down failed ", err)
	}
	return nil
}

// GetClusterIP 获取集群内网 VPN IP
func GetClusterIP() string {
	d, err := device.GetDevice()
	utils.CheckErr(err)
	switch d.Vpn {
	case device.VPNModeOPENVPN:
		tun0interface, err := network.GetMyIP(OpenVPNInterface)
		if err != nil {
			log.WithFields(log.Fields{"Operating": "GetClusterIP"}).Error(err)
			return ""
		}
		return tun0interface
	case device.VPNModeWG:
		wg0interface, err := network.GetMyIP(WireGuardInterface)
		if err != nil {
			log.WithFields(log.Fields{"Operating": "GetClusterIP"}).Error(err)
			return ""
		}
		return wg0interface
	}
	return ""
}
