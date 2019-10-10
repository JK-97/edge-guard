package core

import (
    "jxcore/config/yaml"
    "jxcore/core/device"
    "jxcore/core/hearbeat"
    "jxcore/core/register"
    "jxcore/log"
    "jxcore/lowapi/network"
    "jxcore/lowapi/utils"
    "jxcore/management/updatemanage"
    "jxcore/monitor/dnsdetector"
    "sync"
    "time"
)

var DnsOnce sync.Once

//control the base version 
func BaseCore() {
    //UpdateCore(10)
    startupProgram, err := yaml.LoadYaml(YamlComponentSetting)
    utils.CheckErr(err)
    yaml.ParseAndCheck(startupProgram, "")

}

//control the base version 
func ProCore() {
    var err error
    var mymasterip string
    currentedvice, err := device.GetDevice()
    utils.CheckErr(err)
    for {
        for {
            register.FindMasterFromDHCPServer(currentedvice.WorkerID, currentedvice.Key)
            mymasterip, err = register.GetMyMaster(currentedvice.WorkerID, currentedvice.Key)
            if err == nil {
                break
            }
            log.Error("Register Worker Net", err)
            time.Sleep(3 * time.Second)
        }
        time.Sleep(3 * time.Second)
        go DnsOnce.Do(dnsdetector.RunDnsDetector)
        // VPN 就绪之后 启动 component 按照配置启动(同步工具集合)
        hearbeat.AliveReport(mymasterip)
    }
}

//contrl the update 
func UpdateCore(timeout int) {
    if network.CheckNetwork() {
        starttime := time.Now()
        updateprocess := updatemanage.GetUpdateProcess()
        //updateprocess.UploadVersion()
        pkgneedupdate := updateprocess.CheckUpdate()
        if len(pkgneedupdate) != 0 {
            updateprocess.UpdateSource()
            updateprocess.UpdateComponent(pkgneedupdate)
        }
        for {
            if updateprocess.GetStatus() == updatemanage.FINISHED {
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
