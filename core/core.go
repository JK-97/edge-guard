package core

import (
	"jxcore/config/yaml"
	"jxcore/core/device"
	"jxcore/core/hearbeat"
	"jxcore/core/register"
	log "jxcore/go-utils/logger"
	"jxcore/lowapi/dns"
	"jxcore/lowapi/network"
	"jxcore/lowapi/utils"
	"jxcore/management/updatemanage"
	"jxcore/monitor/dnsdetector"
	"time"
)

func NewJxCore() *JxCore {
	return &JxCore{}
}

func GetJxCore() *JxCore {

	lock.Lock()
	defer lock.Unlock()
	if jxcore == nil {
		jxcore = NewJxCore()
		return jxcore
	}
	return jxcore
}

//control the base version
func (j *JxCore) BaseCore() {
	//UpdateCore(10)
	startupProgram, err := yaml.LoadYaml(YamlComponentSetting)
	utils.CheckErr(err)
	yaml.ParseAndCheck(startupProgram, "")

}

//control the base version
func (j *JxCore) ProCore() {
	var err error
	var mymasterip string
	currentedvice, err := device.GetDevice()
	utils.CheckErr(err)

	SetUp()
	if usbNetworkReachable {
		done := make(chan struct{})
		go linkSubscribe(done)
	}

	for {
		dns.CheckResolvFile()
		go DnsOnce.Do(dnsdetector.RunDnsDetector)
		dns.ResetHostFile(network.GetEthIP())
		for {

			register.FindMasterFromDHCPServer(currentedvice.WorkerID, currentedvice.Key)
			//获取vpn key，连接vpn
			mymasterip, err = register.GetMyMaster(currentedvice.WorkerID, currentedvice.Key)
			//校验新的master是否协力hossts文件
			if err == nil {
				break
			}
			log.Error("Register Worker Net", err)
			time.Sleep(3 * time.Second)
		}
		time.Sleep(3 * time.Second)

		// VPN 就绪之后 启动 component 按照配置启动(同步工具集合)
		hearbeat.AliveReport(mymasterip)
	}
}

//contrl the update
func (j JxCore) UpdateCore(timeout int) {

	for !network.CheckMasterConnect() {
		time.Sleep(3 * time.Second)
		log.Info("Waiting for master connect")
	}
	log.Info("Master Connect")
	if dns.CheckDnsmasqConf() {
		log.Info("Normal Dnsmasq configuration ")
	} else {
		log.Error("Error Dnsmasq configuration ")
	}
	updatemanage.AddAptKey()
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
