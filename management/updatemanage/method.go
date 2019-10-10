package updatemanage

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "jxcore/core/device"
    "jxcore/log"
    "strings"

    "net/http"
    "os/exec"
)

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
    //log.Info(targetinfo.Target)
    return &UpgradeProcess{
        //Target:     targetinfo.Target["target"],
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
    up.ChangeToUpdateSource()
    exec.Command("apt", "updatemanage").Run()
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
    log.Info("now version", up.NowVersion)
    log.Info("target version", up.Target)
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
    if err != nil {
        log.Error(err)
    }
    _, err = http.Post(UPLOADURL, "application/json", bytes.NewReader(respdata))
    if err != nil {
        log.Error(err)
    }

}

func (up *UpgradeProcess) UpdateComponent(componenttoupdate map[string]string) {
    up.ChangeToUpdating()
    for pkgname, pkgversion := range componenttoupdate {
        exec.Command("apt", "autoremove", pkgname).Run()
        log.Info("updating " + pkgname)
        pkginfo := pkgname + "=" + pkgversion
        exec.Command("apt", "install", "-y", pkginfo).Run()
    }
    up.FlushVersionInfo()
    up.ChangeToFinish()
}

func (up *UpgradeProcess) ChangeToFinish() {
    up.Status = FINISHED

}
func (up *UpgradeProcess) ChangeToUpdating() {
    up.Status = UPDATING

}
func (up *UpgradeProcess) ChangeToUpdateSource() {
    up.Status = UPDATESOURCE

}
func (up *UpgradeProcess) SetNewTarget(indentdata []byte) {
    ioutil.WriteFile(TARGETVERSION, indentdata, 0644)
    up.FlushVersionInfo()
}
