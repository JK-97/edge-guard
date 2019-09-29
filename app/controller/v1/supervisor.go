package v1

import (
	"jxcore/app/model/supervisor"
	"jxcore/config"
	"net/http"
)


var s = supervisor.NewSupervisorRPC(config.InterSettings.Supervisor.Host)

//SupervisorAllProcessGET SupervisorAllProcessGET ssss
func SupervisorAllProcessGET(w http.ResponseWriter, r *http.Request) {
	procs, _ := s.GetAllProcessInfo()
	//if err != nil {
	//	panic(err)
	//}

	respondSuccessJSON(procs, w, r)
}

//SupervisorStopProcessSupervisorStopProcess
func SupervisorRestoreProcessPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.StopAllProcesses(true)
		s.ReloadConfig()
		s.StartAllProcesses(false)

		respondSuccessJSON(nil, w, r)
	}
}
