package updateM

import (
    "bytes"
    "encoding/json"
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "jxcore/core/device"
    "jxcore/log"
    "jxcore/lowapi/utils"

    "net/http"
    "os"
    "os/exec"
    "path/filepath"
)

func ParseVersionFile() (versioninfo map[string]string) {
    err := filepath.Walk("/edge", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
            return err
        }
        if !info.IsDir() {
            if len(path) > 7 {

                if path[len(path)-7:] == "version" {

                    fileRawData, err := ioutil.ReadFile(path)
                    if err != nil {
                        log.Error(err)
                    }
                    conponentInfo := ComponentInfo{}
                    yaml.Unmarshal(fileRawData, &conponentInfo)
                    versioninfo[conponentInfo.Name] = conponentInfo.Version
                }
            }

        }
        return err
    })
    if err != nil {
        log.Error(err)
    }
    return versioninfo
}
func NewUpdateProcess() *UpgradeProcess {

    targetdata, err := ioutil.ReadFile("/edge/target")
    if err != nil {
        log.Error(err)
    }
    targetinfo := targetversionfile{}
    json.Unmarshal(targetdata, &targetinfo)
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
    }
    return process
}

func (up *UpgradeProcess) UpdateSource() {
    up.ChangeToUpdateSource()
    exec.Command("apt", "updateM").Run()
}

func (up *UpgradeProcess) GetStatus() UpgradeStatus {
    return up.Status
}

func (up *UpgradeProcess) CheckUpdate() map[string]string {
    var pkgNeedUpdate = make(map[string]string)
    for pkgName, version := range up.Target {
        if up.NowVersion[pkgName] != version {
            pkgNeedUpdate[pkgName] = version
        }
    }
    return pkgNeedUpdate

}

func (up *UpgradeProcess) UploadVersion() {
    url := "http://masterip/api/v1/worker/version_info"
    deviceInfo, err := device.GetDevice()
    utils.CheckErr(err)
    respRawData := Respdatastruct{
        Status:   up.GetStatus().String(),
        WorkerId: deviceInfo.WorkID,
        PkgInfo:  ParseVersionFile(),
    }
    respData, err := json.Marshal(respRawData)
    if err != nil {
        log.Error(err)
    }
    http.Post(url, "application/json", bytes.NewReader(respData))

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
func (up *UpgradeProcess) UpdateComponent(componenttoupdate map[string]string) {
    up.ChangeToUpdating()
    for pkgName, pkgVersion := range componenttoupdate {
        pkgInfo := pkgName + "==" + pkgVersion
        exec.Command("apt", "install", pkgInfo)
    }
    up.ChangeToFinish()
}
