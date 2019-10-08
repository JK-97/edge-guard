package core

import (
    "jxcore/config/yaml"
    "jxcore/journal"
    "jxcore/log"
    "jxcore/lowapi/network"
    "jxcore/lowapi/utils"
    "jxcore/management/updateM"
    "jxcore/monitor/dnsdetector"
    "os"
    "time"
)

//control the base version 
func BaseCore() {

    UpdateCore(30)
    startupProgram, err := yaml.LoadYaml(YamlComponentSetting)
    utils.CheckErr(err)
    yaml.ParseAndCheck(startupProgram, "")

}

//control the base version 
func ProCore() {
    // (pro版本, jxcore hardcode) 启动 VPN, 离线模式下一直重试

    
    dnsdetector.RunDnsDetector()
    // pro板 需要检查 同步工具完整性(按照启动的配置)

    // VPN 就绪之后 启动 component 按照配置启动(同步工具集合)
}

func DnsCore() {

}

//contrl the update 
func UpdateCore(timeout int) {
    if network.CheckNetwork() {
        starttime := time.Now()
        updateprocess := updateM.GetUpdateProcess()
        updateprocess.UploadVersion()
        pkgneedupdate := updateprocess.CheckUpdate()
        if len(pkgneedupdate) != 0 {
            updateprocess.UpdateSource()
            updateprocess.UpdateComponent(pkgneedupdate)
            log.Info("updating")
        }
        for {
            if updateprocess.GetStatus() == updateM.FINISHED {
                break
            }
            if time.Now().Unix() > starttime.Add(time.Duration(timeout)*time.Second).Unix() {
                log.Error("update time out ")
                break
            }
        }
        updateprocess.UploadVersion()
    } else {
        log.Warn("The network is not working properly and automatically enters offline mode.")
    }

}

func CollectJournal(workerID string) {

    ttl := time.Hour * 24 * 30 // 日志只保留 30 天
    journalConfig := map[string]interface{}{
        "rotate-directory": []string{},
    }

    arcFolder := "/data/edgebox/local/logs"
    metaFolder := "/data/edgebox/remote/logs/" + workerID

    os.MkdirAll(arcFolder, 0755)
    os.MkdirAll(metaFolder, 0755)
    journal.RunForever(&journalConfig, 20*time.Minute, arcFolder, metaFolder, ttl)
}
