package updatemanage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jxcore/core/device"

	"jxcore/internal/network/dns"
	"jxcore/lowapi/logger"
	log "jxcore/lowapi/logger"
	"jxcore/lowapi/system"
	"jxcore/lowapi/utils"

	"net/http"
	"os/exec"
	"strings"
)

func AddAptKey() {
	file, err := ioutil.ReadFile(SourceList)

	utils.CheckErr(err)
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.Contains(line, "deb [arch=arm64] http://master.iotedge/public stable main") {
			return
		}
	}
	err = exec.Command("/bin/bash", "-c", "curl http://master.iotedge/public/gpg | sudo apt-key add -").Run()
	utils.CheckErr(err)
	lines = append(lines, "deb [arch=arm64] http://master.iotedge/public stable main")
	out := strings.Join(lines, "\n")
	ioutil.WriteFile(SourceList, []byte(out), 0666)
}

func ParseVersionFile() (versioninfo map[string]string) {
	versionRawInfo, err := ioutil.ReadFile(EDGEVERSIONFILE)
	if err != nil {
		log.Error(err)
	}
	versioninfo = map[string]string{}
	jxtoolsetversion := strings.TrimSpace(string(versionRawInfo))
	versioninfo["jx-toolset"] = jxtoolsetversion
	return versioninfo
}

func NewUpdateProcess() *UpgradeProcess {
	targetdata, err := ioutil.ReadFile(TARGETVERSION)
	if err != nil {
		log.Error(err)
	}
	targetinfo := targetversionfile{}
	json.Unmarshal(targetdata, &targetinfo.Target)

	return &UpgradeProcess{
		Target:     targetinfo.Target,
		NowVersion: ParseVersionFile(),
		Status:     FINISHED,
	}
}

func GetUpdateProcess() *UpgradeProcess {
	lock.Lock()
	defer lock.Unlock()
	if process == nil {
		process = NewUpdateProcess()
		return process
	}
	new := NewUpdateProcess()
	process.Target = new.Target
	process.NowVersion = new.NowVersion
	return process
}

func (up *UpgradeProcess) UpdateSource() {
	up.ChangeToUpdating()
	log.WithFields(log.Fields{"Operating": "Updating"}).Info("Updating Source")
	err := system.RunCommand("apt update")
	if err != nil {
		logger.Warn(err)
	}
}

func (up *UpgradeProcess) GetStatus() UpgradeStatus {
	return up.Status
}

func (up *UpgradeProcess) FlushVersionInfo() {
	up.NowVersion = ParseVersionFile()
}

func (up *UpgradeProcess) FlushTargetVersion() {
	targetdata, err := ioutil.ReadFile(TARGETVERSION)
	if err != nil {
		log.Error(err)
	}
	targetinfo := targetversionfile{}
	json.Unmarshal(targetdata, &targetinfo.Target)
	up.Target = targetinfo.Target
}

func (up *UpgradeProcess) CheckUpdate() map[string]string {
	var pkgneeddate = make(map[string]string)
	log.WithFields(log.Fields{"Operating": "Updating"}).Info("Current Version : ", up.NowVersion)
	log.WithFields(log.Fields{"Operating": "Updating"}).Info("Target Version : ", up.Target)
	for pkgnamme, version := range up.Target {
		if up.NowVersion[pkgnamme] != version {
			pkgneeddate[pkgnamme] = version
		}
	}
	return pkgneeddate
}

func (up *UpgradeProcess) ReportVersion() {
	deviceinfo, _ := device.GetDevice()

	resprawinfo := Respdatastruct{
		Status:   up.GetStatus().String(),
		WorkerId: deviceinfo.WorkerID,
		PkgInfo:  ParseVersionFile(),
	}
	respdata, err := json.Marshal(resprawinfo)
	//log.Info(string(respdata))
	if err != nil {
		log.Error(err)
	}

	ip, port := dns.ParseIpInTxt(UPLOADDOMAIN)
	_, err = http.Post("http://"+ip+":"+port+UPLOADPATH, "application/json", bytes.NewReader(respdata))

	if err != nil {
		log.Error(err)
	}
}

func (up *UpgradeProcess) UpdateComponent(componenttoupdate map[string]string) error {
	for pkgname, pkgversion := range componenttoupdate {
		err := aptInstall(pkgname, pkgversion)
		if err != nil {
			return err
		}
	}
	up.FlushVersionInfo()
	up.ChangeToFinish()
	return nil
}

func (up *UpgradeProcess) ChangeToFinish() {
	up.Status = FINISHED

}

func (up *UpgradeProcess) ChangeToUpdating() {
	up.Status = UPDATING
}

func (up *UpgradeProcess) SetNewTarget(indentdata []byte) {
	ioutil.WriteFile(TARGETVERSION, indentdata, 0644)
	up.FlushVersionInfo()
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
