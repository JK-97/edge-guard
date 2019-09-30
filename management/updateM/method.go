package updateM

import (
    "bytes"
    "encoding/json"
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "jxcore/log"

    "net/http"
    "os"
    "os/exec"
    "path/filepath"
)

func PraseVersionFile()(versioninfo map[string]string) {
    err := filepath.Walk("/edge", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
            return err
        }
        if !info.IsDir() {
            if len(path) > 7 {

                if (path[len(path)-7:] == "version") {

                    filebyte, err := ioutil.ReadFile(path)
                    if err != nil {
                        log.Error(err)
                    }
                    componentinfo := ComponentInfo{}
                    yaml.Unmarshal(filebyte, &componentinfo)
                    versioninfo[componentinfo.Name] = componentinfo.Version
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
        NowVersion: PraseVersionFile(),
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

func (up *UpgradeProcess) GetStatus() UpgradStatus {
    return up.Status
}

func (up *UpgradeProcess) CheckUpdate() map[string]string {
    var pkgneeddate = make(map[string]string)
    for pkgnamme, version := range up.Target {
        if up.NowVersion[pkgnamme] != version {
            pkgneeddate[pkgnamme] = version
        }
    }
    return pkgneeddate

}

func (up *UpgradeProcess) UploadVersion() {
    url := "http://masterip/api/v1/worker/version_info"
    resprawinfo := Respdatastruct{
        Status:   up.GetStatus().String(),
        WorkerId: regeister.DeviceKey.WorkID,
        PkgInfo:  PraseVersionFile(),
    }
    respdata, err := json.Marshal(resprawinfo)
    if err != nil {
        log.Error(err)
    }
    http.Post(url, "application/json", bytes.NewReader(respdata))

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
    for pkgname, pkgversion := range componenttoupdate {
        pkginfo := pkgname + "==" + pkgversion
        exec.Command("apt", "install", pkginfo)
    }
    up.ChangeToFinish()
}
