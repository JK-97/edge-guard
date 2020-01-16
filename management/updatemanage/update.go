package updatemanage

import (
	"jxcore/lowapi/logger"
	"jxcore/lowapi/system"
	"jxcore/lowapi/utils"
	"time"
)

type UpdateStatus string

const (
	FINISHED UpdateStatus = "finished"
	UPDATING UpdateStatus = "updating"

	EDGEVERSIONFILE string = "/edge/VERSION"
	TARGETVERSION   string = "/etc/edgetarget"
	UPLOADDOMAIN    string = "port30111.version-control.ffffffffffffffffffffffff.master.iotedge"
	UPLOADPATH      string = "/api/v1/worker_version"
	SourceList      string = "/etc/apt/sources.list"
)

type VersionInfo map[string]string

type updateManager struct {
	status         UpdateStatus
	targetVersion  VersionInfo
	currentVersion VersionInfo

	targetVersionUpdated chan bool
}

var process *updateManager

func (up *updateManager) updateLoop() {
	for {
		if err := up.tryUpdate(); err != nil {
			logger.Error(err)
		}

		// wait for target version update
		<-up.targetVersionUpdated
	}
}

func (up *updateManager) tryUpdate() error {
	logger.Info("================Checking JxToolset Update===================")

	needUpdate := getNeedUpdate(up.currentVersion, up.targetVersion)
	if len(needUpdate) == 0 {
		logger.Info("No update required")
		return nil
	}

	// TODO: add lock
	if up.status == UPDATING {
		return ErrUpdating
	}
	up.status = UPDATING
	defer func() { up.status = FINISHED }()

	logger.Info("update required")
	utils.RunAndLogError(addAptKey)
	utils.RunAndLogError(updateSource)
	err := up.updateComponents(needUpdate)
	if err != nil {
		logger.Error(err)
	}

	up.currentVersion = parseCurrentVersionFile()
	utils.RunAndLogError(up.ReportVersion)

	if len(getNeedUpdate(up.currentVersion, up.targetVersion)) == 0 {
		//更新成功后才重启
		logger.Info(" update success ")
		go system.RestartJxcoreAfter(5 * time.Second)
	} else {
		logger.Info(" update fail ")
	}
	return nil
}

// 更新组件
func (up *updateManager) updateComponents(needUpdate VersionInfo) error {
	for pkgname, pkgversion := range needUpdate {
		err := aptInstall(pkgname, pkgversion)
		if err != nil {
			logger.Error(err)
		}
	}
	return nil
}
