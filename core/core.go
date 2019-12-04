package core

import (
	"context"
	"io/ioutil"
	"jxcore/config/yaml"
	"jxcore/core/register"
	"jxcore/internal/network"
	"jxcore/internal/network/dns"
	"jxcore/internal/network/iface"
	"jxcore/lowapi/docker"
	"jxcore/lowapi/utils"
	"jxcore/management/updatemanage"
	"os"

	"jxcore/lowapi/logger"
	log "jxcore/lowapi/logger"

	"golang.org/x/sync/errgroup"
)

func ConfigSupervisor() {
	startupProgram := yaml.Config
	yaml.ParseAndCheck(*startupProgram, "")
}

func ConfigNetwork() {
	// NetworkManager 和 systemd-resolved 会更改 /etc/resolv.conf，使dns不可控。需要停止
	network.DisableNetworkManager()
	network.DisableSystemdResolved()

	// nano: /etc/dhcp/dhclient-enter-hooks.d 里的脚本需要被移除。
	// 如果脚本存在，`ifup`时不会更新 /etc/resolv.conf
	dns.RemoveDHCPEnterHooks()

	// 将dhcp 获取的resolv 信息重定向到 /edge/resolv.d
	dns.ApplyDHCPResolveUpdateHooks()
	dns.ResetResolv()
	dns.ResetDNSMasqConf()
	dns.AppendHostnameHosts()

	err := iface.InitIFace()
	if err != nil {
		log.Error(err)
	}
	err = dns.ResetHostFile()
	if err != nil {
		log.Error(err)
	}

	err = docker.EnsureDockerDNSConfig()
	if err != nil {
		log.Error("Failed to configure docker DNS: ", err)
	}
}

func MaintainNetwork(ctx context.Context) error {
	errGroup := errgroup.Group{}
	errGroup.Go(func() error { return iface.MaintainBestIFace(ctx) })
	errGroup.Go(func() error { return register.MaintainMasterConnection(ctx, UpdateCore) })
	return errGroup.Wait()
}

func UpdateCore() {
	logger.Info("================Checking JxToolset Update===================")
	updatemanage.AddAptKey()
	updateprocess := updatemanage.GetUpdateProcess()
	pkgneedupdate := updateprocess.CheckUpdate()
	if len(pkgneedupdate) != 0 {
		updateprocess.UpdateSource()
		updateprocess.UpdateComponent(pkgneedupdate)
	}
	updateprocess.ReportVersion()
}

// applySyncTools 更新配置同步工具
func applySyncTools() {
	if utils.FileExists("/edge/synctools.zip") {
		data, err := ioutil.ReadFile("/edge/synctools.zip")
		if err != nil {
			log.Error(err)
		} else {
			err = utils.Unzip(data, "/edge/mnt")
			if err == nil {
				log.Info("has find the synctools.zip")
				os.Remove("/edge/synctools.zip.old")
				if err = os.Rename("/edge/synctools.zip", "/edge/synctools.zip.old"); err != nil {
					log.Error("Fail to move /edge/synctools.zip to /edge/synctools.zip.old", err)
				}
			}
		}
	}
}
