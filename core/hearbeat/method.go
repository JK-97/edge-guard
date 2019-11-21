package hearbeat

import (
	"fmt"
	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
	"jxcore/core/device"
	"jxcore/internal/network"
	"jxcore/lowapi/utils"
	"net"
	"time"
)

// AliveReport 上报心跳
func AliveReport(masterip string) {
	var msg string
	currentdevice, err := device.GetDevice()
	utils.CheckErr(err)
	vpn, err := currentdevice.Vpn.Interface()
	utils.CheckErr(err)
	logger := log.WithFields(log.Fields{"Operating": "AliveReport"})
	logger.Info("Get VPN IP")
	vpninterface, err := network.GetMyIP(vpn)
	utils.CheckErr(err)
	msg = fmt.Sprintf(vpninterface + ":" + currentdevice.WorkerID)
	for index := 0; index < 5; {
		if err := tryToSend(masterip, msg); err != nil {
			index++
		} else {
			index = 0
		}
	}
	log.WithFields(log.Fields{"Operating": "AliveReport"}).Error(" will get a new master in 5 second ", err)
	time.Sleep(5 * time.Second)
}

func tryToSend(masterip string, msg string) error {
	conn, err := net.Dial("tcp", masterip+":30431")

	defer conn.Close()
	if err != nil {
		log.Error("disconnect my master ", err)
	} else {
		ticker := time.NewTicker(time.Second * HeartBeatInterva)
		defer ticker.Stop()
		for range ticker.C {
			_, err = conn.Write([]byte(msg))
			if err != nil {
				log.WithFields(log.Fields{"Operating": "AliveReport"}).Error("disconnect my master", err)
				//心跳断联,获取新master
				time.Sleep(3 * time.Second)
				break
			}
		}
	}
	return err
}
