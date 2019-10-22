package updatemanage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"jxcore/core/device"

	log "jxcore/go-utils/logger"
	"jxcore/lowapi/dns"
	"jxcore/lowapi/utils"

	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func AddAptKey() {
	//ip, _ := dns.ParseIpInTxt(UPLOADDOMAIN)
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
	exec.Command("apt", "update").Run()
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

func (up *UpgradeProcess) UploadVersion() {

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

func (up *UpgradeProcess) UpdateComponent(componenttoupdate map[string]string) {

	for pkgname, pkgversion := range componenttoupdate {
		pkginfo := pkgname + "=" + pkgversion
		log.WithFields(log.Fields{"Operating": "Updating"}).Info("Installing : ", pkginfo)
		err := exec.Command("/bin/bash", "-c", "aptitude install -o Aptitude::ProblemResolver::SolutionCost='100*canceled-actions,200*removals' -y "+pkginfo).Run()
		utils.CheckErr(err)

	}
	up.FlushVersionInfo()
	up.ChangeToFinish()
	log.Info("restart in 5 second")
	time.Sleep(5 * time.Second)

	syscall.Exit(0)
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
