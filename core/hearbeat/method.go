package hearbeat

import (
    "fmt"
    "jxcore/core/device"
    "jxcore/log"
    "jxcore/lowapi/network"
    "jxcore/lowapi/utils"
    "net"
    "time"
)

// AliveReport 上报心跳
func AliveReport(masterip string)  {
    var msg string
    currentdevice, err := device.GetDevice()
    utils.CheckErr(err)
    vpninterface ,err:= network.GetMyIP(currentdevice.Vpn)
    utils.CheckErr(err)
    msg = fmt.Sprintf(vpninterface + ":" + currentdevice.WorkerID)
    conn, err := net.Dial("tcp", masterip+":30431")
    if err != nil {
        log.WithFields(log.Fields{"Operating": "AliveReport"}).Error("disconnect my master", err)

    }
    log.Info("connected")

    defer conn.Close()
    ticker := time.NewTicker(time.Millisecond * HeartBeatInterva)
    select {
    case <-ticker.C:

        for index := 0; index < 5; index++ {
            for range ticker.C {
                _, err := conn.Write([]byte(msg))
                if err != nil {
                    log.WithFields(log.Fields{"Operating": "AliveReport"}).Error("disconnect my master", err)
                    //心跳断联,获取新master
                    time.Sleep(3 * time.Second)
                    ticker.Stop()
                    break
                }
            }

        }
        log.WithFields(log.Fields{"Operating": "AliveReport"}).Error(" will get a new master in 5 second ", err)
    }
    time.Sleep(5 * time.Second)

}
