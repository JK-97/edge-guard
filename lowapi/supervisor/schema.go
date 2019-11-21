package supervisor

import "net/http"

type emptyArg struct {
}

var empty = new(emptyArg)

type arg1 struct {
	Arg interface{}
}

// SupervisorRPC SupervisorRPC
type SupervisorRPC struct {
	client *http.Client
}

// ProcessInfoResult ProcessInfoResult
type ProcessInfoResult struct {
	Process ProcessInfo
}

// AllProcessInfoResult ProcessInfoResult
type AllProcessInfoResult struct {
	Processes []ProcessInfo
}

// ProcessInfo Get info about a process named name
type ProcessInfo struct {
	Name          string `json:"name" xml:"name"`
	Group         string `json:"group" xml:"group"`
	Description   string `json:"description" xml:"description"`
	Start         int    `json:"start" xml:"start"`
	Stop          int    `json:"stop" xml:"stop"`
	Now           int    `json:"now" xml:"now"`
	State         int    `json:"state" xml:"state"`
	Statename     string `json:"statename" xml:"statename"`
	Spawnerr      string `json:"spawnerr" xml:"spawnerr"`
	Exitstatus    int    `json:"exitstatus" xml:"exitstatus"`
	Logfile       string `json:"logfile" xml:"logfile"`
	StdoutLogfile string `json:"stdout_logfile" xml:"stdout_logfile"`
	StderrLogfile string `json:"stderr_logfile" xml:"stderr_logfile"`
	Pid           int    `json:"pid" xml:"pid"`
}
