package monitor

import (
	"jxcore/component/process"
	"jxcore/config"
	"jxcore/log"
	"jxcore/regeister"
	"os"
	"os/exec"
)

type GateWay struct {
	UpComponent   chan string
	DownComponent chan string
	GwCMD         *exec.Cmd
}

// 组件状态
const (
	UP   = "up"
	DOWN = "down"
	EXIT = "exit"
)

var GW = GateWay{UpComponent: make(chan string, 1), DownComponent: make(chan string, 1)}
var gatewayprocess *process.Process
var gwinitfinished = make(chan bool, 1)

// GWEmitter start the gw
func GWEmitter() *process.Process {
	// exec.Command("dhclient").Run()

	//eth0
	netaddr := regeister.GetEthIP()

	f, err := os.OpenFile(regeister.HostsFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	f.WriteString(netaddr + " edgegw.localhost\n")
	f.WriteString(netaddr + " edgegw.iotedge\n")
	exec.Command("service", "dnsmasq", "restart").Run()

	// f.WriteString(res[0] + " " + "edgegw.localhost")
	f.Close()
	needStopGateway = true
	//gatewaycmd := exec.Command(config.InterSettings.Gateway.Name, "-c", config.InterSettings.Gateway.Config)
	gatewayprocess = process.NewProcess("1", config.InterSettings.Gateway.Name)
	//err = gatewaycmd.Start()
	if err != nil {
		log.Error(err)
	}
	//GW.GwCMD = gatewaycmd

	log.Error("********************Dividing line*********************")
	gatewayprocess.Start(true)
	GW.GwCMD = gatewayprocess.GetCmd()
	//log.WithFields(log.Fields{"program": "gateway"}).Info("GateWay start pid:", GW.GwCMD.Process.Pid)
	gwinitfinished <- true
	return gatewayprocess
}

var needStopGateway bool

// StopGateway 关闭 gateway
func StopGateway() {
	log.Info("Begin Exit GateWayMonitor")
	needStopGateway = true
	GW.DownComponent <- EXIT
	<-Finished
	gatewayprocess.Stop(true)
	log.Info("Finish Exit GateWayMonitor")
}

// GateWayMonitor 监控 GateWay 进程
func GateWayMonitor() {
	<-gwinitfinished
	for {
		//for {
		//	log.Error(gatewayprocess.GetState().String())
		//	if gatewayprocess.GetState() != process.RUNNING || gatewayprocess.GetState()!= process.STARTING  {
		//		break
		//	}
		//}
		GW.UpComponent <- UP
		GW.GwCMD.Process.Wait()

		log.WithFields(log.Fields{"program": "gateway"}).Error("GateWay finish pid:", GW.GwCMD.Process.Pid)
		log.WithFields(log.Fields{"program": "gateway"}).Info("gateway will auto start at 5 secode")
		log.Error("down")
		GW.DownComponent <- DOWN
		<-Finished

		for {
			gatewayprocess.Wait()

			if gatewayprocess.GetState() == process.RUNNING {
				GW.GwCMD = gatewayprocess.GetCmd()
				break
			}
		}
	}

	//GWEmitter()
}
