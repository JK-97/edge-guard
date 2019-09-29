package monitor

import (
	"jxcore/component"
	"jxcore/component/process"
	"jxcore/log"
	"jxcore/regeister"
)

//ComponentMonitor FROM settings
type CMD struct {
	args []string
}

var Finished = make(chan bool, 0)

func ComponentMonitor() {
	<-regeister.Connectable
	component.ComponentPidInfo.Gpid = make([]*process.Process, 0)
	for {
		select {
		case msg1 := <-GW.DownComponent:
			log.Error(msg1, " Component....")
			component.StopComponent("")
			// for _, perprocess := range component.ComponentPidInfo.Gpid {
			// 	perprocess.Stop(true)
			// }
			log.Error(msg1, " Component", " finish")
			component.ComponentPidInfo.Gpid = make([]*process.Process, 0)
			Finished <- true
		case msg2 := <-GW.UpComponent:
			log.Info(msg2, " Component....")
			component.ComponentEmiter()
			log.Info(msg2, " Component", " finish")

			//Finished<-true
		}

	}

}
