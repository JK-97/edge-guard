package core

import (
	"context"
	"io/ioutil"
	"jxcore/config/yaml"
	"jxcore/core/register"
	"jxcore/internal/network/dns"
	"jxcore/internal/network/iface"
	"jxcore/lowapi/system"
	"jxcore/lowapi/utils"
	"jxcore/management/updatemanage"
	"os"
	"time"

	"jxcore/lowapi/logger"
	log "jxcore/lowapi/logger"

	"golang.org/x/sync/errgroup"
)

func ConfigSupervisor() {
	startupProgram := yaml.Config
	yaml.ParseAndCheck(*startupProgram, "")
}

func MaintainNetwork(ctx context.Context, noUpdate bool) error {
	dns.TrySetupDnsConfig()

	errGroup := errgroup.Group{}

	// 按优先级切换网口
	errGroup.Go(func() error { return iface.MaintainBestIFace(ctx) })

	// 第一次连接master成功，检查固件更新
	onFirstConnect := func() {
		updatemanage.GetUpdateProcess().ReportVersion()
		if !noUpdate {
			CheckCoreUpdate()
		}
	}
	errGroup.Go(func() error { return register.MaintainMasterConnection(ctx, onFirstConnect) })

	return errGroup.Wait()
}

func CheckCoreUpdate() {
	logger.Info("================Checking JxToolset Update===================")
	updatemanage.AddAptKey()
	updateprocess := updatemanage.GetUpdateProcess()
	pkgneedupdate := updateprocess.CheckUpdate()
	if len(pkgneedupdate) != 0 {
		updateprocess.UpdateSource()
		err := updateprocess.UpdateComponent(pkgneedupdate)
		if err != nil {
			logger.Warn(err)
			return
		}
		updateprocess.ReportVersion()
		system.RestartJxcoreAfter(5 * time.Second)
	}
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
