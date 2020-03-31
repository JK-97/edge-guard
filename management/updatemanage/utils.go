package updatemanage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/JK-97/edge-guard/lowapi/logger"
	log "github.com/JK-97/edge-guard/lowapi/logger"

	"os/exec"
	"strings"
)

var GpgKeyInsertError = errors.New("密钥导入失败")

//添加 respository key
func addAptKey() error {
	output, err := exec.Command("/bin/bash", "-c", "curl http://master.iotedge/public/gpg | sudo apt-key add -").Output()
	logger.Info(string(output))
	if err != nil {
		return GpgKeyInsertError
	}
	// 向sourcelist 添加源记录
	file, err := ioutil.ReadFile(SourceList)
	if err != nil {
		return err
	}
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.Contains(line, "deb [arch=arm64] http://master.iotedge/public stable main") {
			return nil
		}
	}
	lines = append(lines, "deb [arch=arm64] http://master.iotedge/public stable main")
	_ = ioutil.WriteFile(SourceList, []byte(strings.Join(lines, "\n")), 0666)
	return nil
}

// 解析当前版本信息文件
func parseCurrentVersionFile() (versioninfo map[string]string) {
	versionRawInfo, err := ioutil.ReadFile(EDGEVERSIONFILE)
	if err != nil {
		logger.Error(err)
	}
	versioninfo = map[string]string{}
	jxtoolsetversion := strings.TrimSpace(string(versionRawInfo))
	versioninfo["jx-toolset"] = jxtoolsetversion
	return versioninfo
}

type targetversionfile struct {
	Target map[string]string `json:"target"`
}

// 解析目标版本信息文件
func parseTargetVersionFile() VersionInfo {
	targetdata, err := ioutil.ReadFile(TARGETVERSION)
	if err != nil {
		logger.Error(err)
	}
	targetinfo := targetversionfile{}
	_ = json.Unmarshal(targetdata, &targetinfo.Target)
	return targetinfo.Target
}

// 更新debian 源
func updateSource() error {
	logger.WithFields(logger.Fields{"Operating": "Updating"}).Info("Updating Source")
	err := exec.Command("apt", "update").Run()
	if err != nil {
		return err
	}
	return nil
}

// 检查需要更新的包
func getNeedUpdate(currentVersion VersionInfo, targetVersion VersionInfo) VersionInfo {
	needUpdate := VersionInfo{}
	logger.WithFields(logger.Fields{"Operating": "Updating"}).Info("Current Version : ", currentVersion)
	logger.WithFields(logger.Fields{"Operating": "Updating"}).Info("Target Version : ", targetVersion)
	for pkgnamme, version := range targetVersion {
		if currentVersion[pkgnamme] != version {
			needUpdate[pkgnamme] = version
		}
	}
	return needUpdate
}

const noPackageInstalledPrompt = "No packages will be installed, upgraded, or removed."

var ErrNoPackageInstalled = fmt.Errorf(noPackageInstalledPrompt)

func aptInstall(pkgname, pkgversion string) error {
	pkginfo := pkgname + "=" + pkgversion
	log.WithFields(log.Fields{"Operating": "Updating"}).Info("Installing : ", pkginfo)

	cmd := "aptitude install -o Aptitude::ProblemResolver::SolutionCost='100*canceled-actions,200*removals' -y " + pkginfo
	output, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	if strings.Contains(string(output), noPackageInstalledPrompt) {
		return ErrNoPackageInstalled
	}
	return nil
}
