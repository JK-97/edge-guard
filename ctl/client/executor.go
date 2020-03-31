package client

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/JK-97/edge-guard/ctl/rpc"

	"github.com/ochinchina/supervisord/types"
)

type RpcExector struct {
	ServerUrl string `short:"s" long:"serverurl" description:"URL on which supervisord server is listening"`
	User      string `short:"u" long:"user" description:"the user name"`
	Password  string `short:"P" long:"password" description:"the password"`
	Verbose   bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
}

type StatusCommand struct {
}

type StartCommand struct {
}

type StopCommand struct {
}

type RestartCommand struct {
}

type ShutdownCommand struct {
}

type ReloadCommand struct {
}

type PidCommand struct {
}

type SignalCommand struct {
}

type TailCommand struct {
}

func NewRpcExcutior(serveurl, user, password string) *RpcExector {
	return &RpcExector{
		ServerUrl: serveurl,
		User:      user,
		Password:  password,
	}
}

func (x *RpcExector) getServerUrl() string {
	return x.ServerUrl
}
func (x *RpcExector) getUser() string {
	return x.User
}
func (x *RpcExector) getPassword() string {

	return x.Password
}

func (x *RpcExector) createRpcClient() *rpc.XmlRPCClient {
	rpcc := rpc.NewXmlRPCClient(x.getServerUrl(), x.Verbose)
	rpcc.SetUser(x.getUser())
	rpcc.SetPassword(x.getPassword())
	return rpcc
}

func (x *RpcExector) Execute(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	} else if s == "quit" || s == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}
	args := strings.Split(s, " ")
	length := len(args)
	rpcc := x.createRpcClient()
	firstword := args[0]

	switch firstword {
	case "status":
		x.status(rpcc, args[1:])
	case "start", "stop":
		if length >= 1 {
			x.startStopProcesses(rpcc, firstword, args[1:])
			return
		}
	case "reset":
		x.shutdown(rpcc)
	case "restart":
		x.restartProcesses(rpcc, args[1:])
	// case "reload":
	// 	x.reload(rpcc)
	case "signal":
		sig_name, processes := args[1], args[2:]
		x.signal(rpcc, sig_name, processes)
	case "tail":
		if length == 3 {
			thirdword := args[2]
			switch thirdword {
			case "stderr", "stdout":
				tailProcessLog(args[1:])
			}
		}
	case "pid":
		if length == 2 {
			x.getPid(rpcc, args[1])
			return
		}
	case "log":
		cmd := exec.Command("tail", "-f", "/edge/logs/edge-guard.log")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
		}
	case "help":

	default:
		fmt.Println("unknown command")
	}
}

// get the status of processes
func (x *RpcExector) status(rpcc *rpc.XmlRPCClient, processes []string) {
	processesMap := make(map[string]bool)
	for _, process := range processes {
		processesMap[process] = true
	}
	if reply, err := rpcc.GetAllProcessInfo(); err == nil {
		x.showProcessInfo(&reply, processesMap)
	} else {
		os.Exit(1)
	}
}

func (x *RpcExector) getAllProcessesName(rpcc *rpc.XmlRPCClient) (processList []string) {
	if reply, err := rpcc.GetAllProcessInfo(); err == nil {
		processList = x.getProcessInfo(&reply)
	} else {
		os.Exit(1)
	}
	return processList
}

// start or stop the processes
// verb must be: start or stop
func (x *RpcExector) startStopProcesses(rpcc *rpc.XmlRPCClient, verb string, processes []string) {
	state := map[string]string{
		"start": "started",
		"stop":  "stopped",
	}
	x._startStopProcesses(rpcc, verb, processes, state[verb], true)
}

func (x *RpcExector) _startStopProcesses(rpcc *rpc.XmlRPCClient, verb string, processes []string, state string, showProcessInfo bool) {
	if len(processes) <= 0 {
		fmt.Printf("Please specify process for %s\n", verb)
	}
	for _, pname := range processes {
		if pname == "all" {
			reply, err := rpcc.ChangeAllProcessState(verb)
			if err == nil {
				if showProcessInfo {
					x.showProcessInfo(&reply, make(map[string]bool))
				}
			} else {
				fmt.Printf("Fail to change all process state to %s", state)
			}
		} else {
			if reply, err := rpcc.ChangeProcessState(verb, pname); err == nil {
				if showProcessInfo {
					fmt.Printf("%s: ", pname)
					if !reply.Value {
						fmt.Printf("not ")
					}
					fmt.Printf("%s\n", state)
				}
			} else {
				fmt.Printf("%s: failed [%v]\n", pname, err)
				os.Exit(1)
			}
		}
	}
}

func (x *RpcExector) restartProcesses(rpcc *rpc.XmlRPCClient, processes []string) {
	x._startStopProcesses(rpcc, "stop", processes, "stopped", false)
	x._startStopProcesses(rpcc, "start", processes, "restarted", true)
}

// shutdown the supervisord
func (x *RpcExector) shutdown(rpcc *rpc.XmlRPCClient) {
	if reply, err := rpcc.Shutdown(); err == nil {
		if reply.Value {
			fmt.Printf("Shut Down\n")
		} else {
			fmt.Printf("Hmmm! Something gone wrong?!\n")
		}
	} else {
		os.Exit(1)
	}
}

// reload all the programs in the supervisord
func (x *RpcExector) reload(rpcc *rpc.XmlRPCClient) {
	if reply, err := rpcc.ReloadConfig(); err == nil {

		if len(reply.AddedGroup) > 0 {
			fmt.Printf("Added Groups: %s\n", strings.Join(reply.AddedGroup, ","))
		}
		if len(reply.ChangedGroup) > 0 {
			fmt.Printf("Changed Groups: %s\n", strings.Join(reply.ChangedGroup, ","))
		}
		if len(reply.RemovedGroup) > 0 {
			fmt.Printf("Removed Groups: %s\n", strings.Join(reply.RemovedGroup, ","))
		}
	} else {
		os.Exit(1)
	}
}

// send signal to one or more processes
func (x *RpcExector) signal(rpcc *rpc.XmlRPCClient, sig_name string, processes []string) {
	for _, process := range processes {
		if process == "all" {
			reply, err := rpcc.SignalAll(process)
			if err == nil {
				x.showProcessInfo(&reply, make(map[string]bool))
			} else {
				fmt.Printf("Fail to send signal %s to all process", sig_name)
				os.Exit(1)
			}
		} else {
			reply, err := rpcc.SignalProcess(sig_name, process)
			if err == nil && reply.Success {
				fmt.Printf("Succeed to send signal %s to process %s\n", sig_name, process)
			} else {
				fmt.Printf("Fail to send signal %s to process %s\n", sig_name, process)
				os.Exit(1)
			}
		}
	}
}

// get the pid of running program
func (x *RpcExector) getPid(rpcc *rpc.XmlRPCClient, process string) {
	procInfo, err := rpcc.GetProcessInfo(process)
	if err != nil {
		fmt.Printf("program '%s' not found\n", process)
		os.Exit(1)
	} else {
		fmt.Printf("%d\n", procInfo.Pid)
	}
}

// check if group name should be displayed
func (x *RpcExector) showGroupName() bool {
	val, ok := os.LookupEnv("SUPERVISOR_GROUP_DISPLAY")
	if !ok {
		return false
	}

	val = strings.ToLower(val)
	return val == "yes" || val == "true" || val == "y" || val == "t" || val == "1"
}

// tail the process stdout log
func tailProcessLog(args []string) {
	process := args[0]
	if strings.Contains(args[0], ":") {
		process = strings.Split(args[0], ":")[0]
	}
	logLevel := "stderr"
	if len(args) == 2 {
		logLevel = args[1]
	}
	cmd := exec.Command("tail", "-f", "/edge/logs/"+process+"_"+logLevel+".log.0")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err.Error())
	}
	return
}

func (x *RpcExector) showLogFlow(processTailLog *rpc.ProcessTailLog) {
	if processTailLog.LogData != "" {
		fmt.Println(processTailLog.LogData)
	}

}
func (x *RpcExector) showProcessInfo(reply *rpc.AllProcessInfoReply, processesMap map[string]bool) {
	for _, pinfo := range reply.Value {
		description := pinfo.Description
		if strings.ToLower(description) == "<string></string>" {
			description = ""
		}
		if x.inProcessMap(&pinfo, processesMap) {
			processName := pinfo.GetFullName()
			if !x.showGroupName() {
				processName = pinfo.Name
			}
			fmt.Printf("%s%-33s%-10s%s%s\n", x.getANSIColor(pinfo.Statename), processName, pinfo.Statename, description, "\x1b[0m")
		}
	}
}
func (x *RpcExector) getProcessInfo(reply *rpc.AllProcessInfoReply) (processesList []string) {
	for _, pinfo := range reply.Value {
		description := pinfo.Description
		if strings.ToLower(description) == "<string></string>" {
			description = ""
		}
		processName := pinfo.GetFullName()
		processesList = append(processesList, processName)
	}
	return processesList
}

func (x *RpcExector) inProcessMap(procInfo *types.ProcessInfo, processesMap map[string]bool) bool {
	if len(processesMap) <= 0 {
		return true
	}
	for procName, _ := range processesMap {
		if procName == procInfo.Name || procName == procInfo.GetFullName() {
			return true
		}

		// check the wildcast '*'
		pos := strings.Index(procName, ":")
		if pos != -1 {
			groupName := procName[0:pos]
			programName := procName[pos+1:]
			if programName == "*" && groupName == procInfo.Group {
				return true
			}
		}
	}
	return false
}

func (x *RpcExector) getANSIColor(statename string) string {
	if statename == "RUNNING" {
		// green
		return "\x1b[0;32m"
	} else if statename == "BACKOFF" || statename == "FATAL" {
		// red
		return "\x1b[0;31m"
	} else {
		// yellow
		return "\x1b[1;33m"
	}
}
