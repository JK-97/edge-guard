package process

import (
	"fmt"
	"io"
	"jxcore/component/events"
	"jxcore/component/logger"
	"jxcore/component/signals"
	"jxcore/log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type ProcessState int

const (
	STOPPED  ProcessState = iota
	STARTING              = 10
	RUNNING               = 20
	BACKOFF               = 30
	STOPPING              = 40
	EXITED                = 100
	FATAL                 = 200
	UNKNOWN               = 1000
)

func (p ProcessState) String() string {
	switch p {
	case STOPPED:
		return "STOPPED"
	case STARTING:
		return "STARTING"
	case RUNNING:
		return "RUNNING"
	case BACKOFF:
		return "BACKOFF"
	case STOPPING:
		return "STOPPING"
	case EXITED:
		return "EXITED"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type Process struct {
	supervisor_id string
	configpath    string
	name          string
	ground        string
	bin           string
	cmd           *exec.Cmd
	startTime     time.Time
	stopTime      time.Time
	state         ProcessState
	//true if process is starting
	inStart    bool
	lock       sync.RWMutex
	stdin      io.WriteCloser
	stopByUser bool
	retryTimes *int32
	StdoutLog  logger.Logger
	StderrLog  logger.Logger
}

func NewProcess(supervisor_id string, binpath string) *Process {
	proc := &Process{
		supervisor_id: supervisor_id,
		configpath:    binpath + ".cfg",
		bin:           binpath,
		ground:        "",
		name:          "",
		cmd:           nil,
		startTime:     time.Unix(0, 0),
		stopTime:      time.Unix(0, 0),
		state:         STOPPED,
		inStart:       false,
		stopByUser:    false,
		retryTimes:    new(int32)}
	//同步工具放置在/edge/mnt目录下

	res := strings.Split(binpath, "/")
	proc.name = res[len(res)-1]
	proc.cmd = nil
	return proc
}

func (p *Process) GetCmd() *exec.Cmd {
	return p.cmd
}

func (p *Process) Wait() {
	p.cmd.Wait()
}

// IsStopByUser 指示程序是否被用户主动关闭
func (p *Process) IsStopByUser() bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.stopByUser
}

func (p *Process) Start(wait bool) {
	//log.WithFields(log.Fields{"programM": p.GetName()}).Info("try to start programM")
	p.lock.Lock()
	if p.stopByUser {
		p.lock.Unlock()
		return
	}

	if p.inStart {
		log.WithFields(log.Fields{"programM": p.GetName()}).Info("Don't start programM again, programM is already started")
		p.lock.Unlock()
		return
	}

	p.inStart = true
	p.stopByUser = false
	p.lock.Unlock()

	var runCond *sync.Cond
	finished := false
	if wait {
		runCond = sync.NewCond(&sync.Mutex{})
		runCond.L.Lock()
	}

	go func() {

		for {
			if wait {
				runCond.L.Lock()
			}
			p.run(func() {
				finished = true
				if wait {
					runCond.L.Unlock()
					runCond.Signal()
				}

			})
			//avoid print too many logs if fail to start programM too quickly
			if time.Now().Unix()-p.startTime.Unix() < 2 {
				time.Sleep(5 * time.Second)
			}
			if p.stopByUser {
				log.WithFields(log.Fields{"programM": p.GetName()}).Info("Stopped by user, don't start it again")
				break
			}

		}
		p.lock.Lock()
		p.inStart = false
		p.lock.Unlock()
	}()

	if wait && !finished {
		runCond.Wait()
		runCond.L.Unlock()
	}
}
func (p *Process) GetName() string {
	return p.name
}
func (p *Process) isRunning() bool {
	if p.cmd != nil && p.cmd.ProcessState != nil {
		if runtime.GOOS == "windows" {
			proc, err := os.FindProcess(p.cmd.Process.Pid)
			return proc != nil && err == nil
		} else {
			return p.cmd.Process.Signal(syscall.Signal(0)) == nil
		}
	}
	return false
}
func (p *Process) getStartSeconds() int64 {
	return int64(1)
}

func (p *Process) createProgramCommand() (err error) {
	if p.name == "telegraf" {
		p.cmd = exec.Command(p.bin, "--config", p.configpath)
	} else if p.name == "mcuserver" {
		p.cmd = exec.Command(p.bin)
	} else if p.name == "watchdog" {
		p.cmd = exec.Command("/bin/bash", p.bin)
	} else if p.name == "powermanagement" {
		p.cmd = exec.Command("/bin/bash", p.bin)
	} else if p.name == "db" {
		p.cmd = exec.Command(p.bin, "serve", "--repo", "mongodb://172.17.0.1:27017", "--src", "mongodb://172.17.0.1:27017")
	} else if p.name == "fs" {
		p.cmd = exec.Command("/bin/bash", p.bin)
	} else {
		p.cmd = exec.Command(p.bin, "-c", p.configpath)
	}
	p.setLog()
	return
}

//
//func (p *Process) handleOutput(stdout io.ReadCloser, stderr io.ReadCloser) error {
//
//	var stdoutFilename string
//	if stdoutFilename == "" {
//		stdoutFilename = fmt.Sprintf("/tmp/%s %s-%d.stdout", p.name, filepath.Base(p.cmd.Path), p.cmd.Process.Pid)
//	}
//
//	var stderrFilename string
//	if stderrFilename == "" {
//		stderrFilename = fmt.Sprintf("/tmp/%s %s-%d.stderr",p.name, filepath.Base(p.cmd.Path), p.cmd.Process.Pid)
//	}
//
//	p..Info(fmt.Sprintf("stdout writing to %s", stdoutFilename))
//	p.StderrLog.Info(fmt.Sprintf("stderr writing to %s", stderrFilename))
//	go pipeToFile(log., stdout, stdoutFilename)
//	go pipeToFile(p.StderrLog, stderr, stderrFilename)
//	return nil
//}
//
//func pipeToFile(lg log.Logger, pipe io.ReadCloser, filename string) {
//	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
//	if err != nil {
//		lg.Error(err.Error())
//		return
//	}
//	defer file.Close()
//
//	if _, err := io.Copy(file, pipe); err != nil {
//		lg.Error(err.Error())
//	}
//}
//
func (p *Process) createStdoutLogEventEmitter() logger.LogEventEmitter {
	//if p.config.GetBytes("stdout_capture_maxbytes", 0) <= 0 && false {
	//	return logger.NewStdoutLogEventEmitter(p.config.GetProgramName(), p.config.GetGroupName(), func() int {
	//		return p.GetPid()
	//	})
	//}
	return logger.NewNullLogEventEmitter()
}

func (p *Process) createStderrLogEventEmitter() logger.LogEventEmitter {
	//if p.config.GetBytes("stderr_capture_maxbytes", 0) <= 0 && p.config.GetBool("stderr_events_enabled", false) {
	//	return logger.NewStdoutLogEventEmitter(p.config.GetProgramName(), p.config.GetGroupName(), func() int {
	//		return p.GetPid()
	//	})
	//}
	return logger.NewNullLogEventEmitter()
}

func (p *Process) GetGroup() string {
	return ""
}
func (p *Process) setLog() {

	p.StdoutLog = p.createLogger(p.GetStdoutLogfile(),
		int64(50*1024*1024),
		10,
		p.createStdoutLogEventEmitter())
	capture_bytes := 1000
	if capture_bytes > 0 {
		//log.WithFields(log.Fields{"programM": p.GetName()}).Info("capture stdout process communication")
		p.StdoutLog = logger.NewLogCaptureLogger(p.StdoutLog,
			capture_bytes,
			"PROCESS_COMMUNICATION_STDOUT",
			p.GetName(),
			p.GetGroup())

		p.cmd.Stdout = p.StdoutLog
		//p.config.GetBool("redirect_stderr", false)
		if false {
			p.StderrLog = p.StdoutLog
		} else {
			p.StderrLog = p.createLogger(p.GetStderrLogfile(),
				int64(50*1024*1024),
				10,
				p.createStderrLogEventEmitter())
		}

		capture_bytes = 5000

		if capture_bytes > 0 {
			//log.WithFields(log.Fields{"programM": p.GetName()}).Info("capture stderr process communication")
			p.StderrLog = logger.NewLogCaptureLogger(p.StdoutLog,
				capture_bytes,
				"PROCESS_COMMUNICATION_STDERR",
				p.GetName(),
				p.GetGroup())
		}

		p.cmd.Stderr = p.StderrLog

	}
	//else if p.config.IsEventListener() {
	//	in, err := p.cmd.StdoutPipe()
	//	if err != nil {
	//		log.WithFields(log.Fields{"eventListener": p.config.GetEventListenerName()}).Error("fail to get stdin")
	//		return
	//	}
	//	out, err := p.cmd.StdinPipe()
	//	if err != nil {
	//		log.WithFields(log.Fields{"eventListener": p.config.GetEventListenerName()}).Error("fail to get stdout")
	//		return
	//	}
	//	events := strings.Split(p.config.GetString("events", ""), ",")
	//	for i, event := range events {
	//		events[i] = strings.TrimSpace(event)
	//	}
	//	p.cmd.Stderr = os.Stderr
	//
	//	p.registerEventListener(p.config.GetEventListenerName(),
	//		events,
	//		in,
	//		out)
	//}
}

const folder = "/edge/logs/"

func (p *Process) createLogger(logFile string, maxBytes int64, backups int, logEventEmitter logger.LogEventEmitter) logger.Logger {
	return logger.NewLogger(p.GetName(), logFile, logger.NewNullLocker(), maxBytes, backups, logEventEmitter)
}

func (p *Process) GetStdoutLogfile() string {
	fileName := folder + p.name + ".stdout_logfile"
	expandFile, err := Path_expand(fileName)
	if err != nil {
		return fileName
	}
	return expandFile
}
func (p *Process) GetStderrLogfile() string {
	fileName := folder + p.name + ".stderr_logfile"
	expandFile, err := Path_expand(fileName)
	if err != nil {
		return fileName
	}
	return expandFile
}

func (p *Process) getRestartPause() int {
	return 0
}
func (p *Process) monitorProgramIsRunning(endTime time.Time, monitorExited *int32, programExited *int32) {
	// if time is not expired
	for time.Now().Before(endTime) && atomic.LoadInt32(programExited) == 0 {
		time.Sleep(time.Duration(800) * time.Millisecond)
	}
	atomic.StoreInt32(monitorExited, 1)

	p.lock.Lock()
	defer p.lock.Unlock()
	// if the programM does not exit
	if atomic.LoadInt32(programExited) == 0 && p.state == STARTING {
		if p.stopByUser {
			return
		}
		log.WithFields(log.Fields{"programM": p.GetName()}).Info("success to start programM at pid: ", p.cmd.Process.Pid)
		p.changeStateTo(RUNNING)
	}
}
func (p *Process) waitForExit(startSecs int64) {
	err := p.cmd.Wait()
	if err != nil {
		//log.WithFields(log.Fields{"programM": p.GetName()}).Warn("fail to wait for programM exit")
	} else if p.cmd.ProcessState != nil {
		log.WithFields(log.Fields{"programM": p.GetName()}).Warnf("programM stopped with status:%v", p.cmd.ProcessState)
	} else {
		log.WithFields(log.Fields{"programM": p.GetName()}).Warn("programM stopped")
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.stopTime = time.Now()
}
func (p *Process) run(finishCb func()) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// check if the programM is in running state
	if p.isRunning() {
		log.WithFields(log.Fields{"programM": p.GetName()}).Info("Don't start programM because it is running")
		finishCb()
		return

	}
	p.startTime = time.Now()
	atomic.StoreInt32(p.retryTimes, 0)
	startSecs := p.getStartSeconds()
	restartPause := p.getRestartPause()
	var once sync.Once

	// finishCb can be only called one time
	finishCbWrapper := func() {
		once.Do(finishCb)
	}
	//process is not expired and not stoped by user
	for !p.stopByUser {
		if restartPause > 0 && atomic.LoadInt32(p.retryTimes) != 0 {
			//pause
			p.lock.Unlock()
			log.WithFields(log.Fields{"programM": p.GetName()}).Info("don't restart the programM, start it after ", restartPause, " seconds")
			time.Sleep(time.Duration(restartPause) * time.Second)
			p.lock.Lock()
		}
		endTime := time.Now().Add(time.Duration(startSecs) * time.Second)
		p.changeStateTo(STARTING)
		atomic.AddInt32(p.retryTimes, 1)

		err := p.createProgramCommand()
		if err != nil {
			p.failToStartProgram("fail to create programM", finishCbWrapper)
			break
		}
		//stdout, err := p.cmd.StderrPipe()
		//p.StdoutLog = stdout
		//lg:=log.NewLogger(config.Config())
		//stdout, _ := p.cmd.StdoutPipe()
		//stderr, _ := p.cmd.StderrPipe()
		err = p.cmd.Start()

		//p.handleOutput(stdout, stderr)

		if err != nil {
			if atomic.LoadInt32(p.retryTimes) >= p.getStartRetries() {
				p.failToStartProgram(fmt.Sprintf("fail to start programM with error:%v", err), finishCbWrapper)
				break
			} else {
				log.WithFields(log.Fields{"programM": p.GetName()}).Info("fail to start programM with error:", err)
				p.changeStateTo(BACKOFF)
				continue
			}
		}

		monitorExited := int32(0)
		programExited := int32(0)
		//Set startsec to 0 to indicate that the programM needn't stay
		//running for any particular amount of time.
		if startSecs <= 0 {
			log.WithFields(log.Fields{"programM": p.GetName()}).Info("success to start programM")
			p.changeStateTo(RUNNING)
			go finishCbWrapper()
		} else {
			go func() {
				p.monitorProgramIsRunning(endTime, &monitorExited, &programExited)
				finishCbWrapper()
			}()
		}
		log.WithFields(log.Fields{"programM": p.GetName()}).Debug("wait programM exit")
		p.lock.Unlock()
		p.waitForExit(startSecs)

		atomic.StoreInt32(&programExited, 1)
		// wait for monitor thread exit
		for atomic.LoadInt32(&monitorExited) == 0 {
			time.Sleep(time.Duration(10) * time.Millisecond)
		}

		p.lock.Lock()

		// if the programM still in running after startSecs
		if p.state == RUNNING {
			p.changeStateTo(EXITED)
			log.WithFields(log.Fields{"programM": p.GetName()}).Warn("programM exited")
			break
		} else {
			p.changeStateTo(BACKOFF)
		}

		// The number of serial failure attempts that supervisord will allow when attempting to
		// start the programM before giving up and putting the process into an FATAL state
		// first start time is not the retry time
		if atomic.LoadInt32(p.retryTimes) >= p.getStartRetries() {
			p.failToStartProgram(fmt.Sprintf("fail to start programM because retry times is greater than %d", p.getStartRetries()), finishCbWrapper)
			break
		}
	}

}
func (p *Process) getStartRetries() int32 {
	return 0
}

func (p *Process) failToStartProgram(reason string, finishCb func()) {
	log.WithFields(log.Fields{"programM": p.GetName()}).Errorf(reason)
	p.changeStateTo(FATAL)
	finishCb()
}

// Get the process state
func (p *Process) GetState() ProcessState {
	return p.state
}
func (p *Process) getExitCodes() []int {
	//strExitCodes := strings.Split(p.config.GetString("exitcodes", "0,2"), ",")
	strExitCodes := []string{"123"}
	result := make([]int, 0)
	for _, val := range strExitCodes {
		i, err := strconv.Atoi(val)
		if err == nil {
			result = append(result, i)
		}
	}
	return result
}

func (p *Process) inExitCodes(exitCode int) bool {
	for _, code := range p.getExitCodes() {
		if code == exitCode {
			return true
		}
	}
	return false
}

func (p *Process) changeStateTo(procState ProcessState) {

	progName := p.name
	groupName := p.ground
	if procState == STARTING {
		events.EmitEvent(events.CreateProcessStartingEvent(progName, groupName, p.state.String(), int(atomic.LoadInt32(p.retryTimes))))
	} else if procState == RUNNING {
		events.EmitEvent(events.CreateProcessRunningEvent(progName, groupName, p.state.String(), p.cmd.Process.Pid))
	} else if procState == BACKOFF {
		events.EmitEvent(events.CreateProcessBackoffEvent(progName, groupName, p.state.String(), int(atomic.LoadInt32(p.retryTimes))))
	} else if procState == STOPPING {
		events.EmitEvent(events.CreateProcessStoppingEvent(progName, groupName, p.state.String(), p.cmd.Process.Pid))
	} else if procState == EXITED {
		exitCode, err := p.getExitCode()
		expected := 0
		if err == nil && p.inExitCodes(exitCode) {
			expected = 1
		}
		events.EmitEvent(events.CreateProcessExitedEvent(progName, groupName, p.state.String(), expected, p.cmd.Process.Pid))
	} else if procState == FATAL {
		events.EmitEvent(events.CreateProcessFatalEvent(progName, groupName, p.state.String()))
	} else if procState == STOPPED {
		events.EmitEvent(events.CreateProcessStoppedEvent(progName, groupName, p.state.String(), p.cmd.Process.Pid))
	} else if procState == UNKNOWN {
		events.EmitEvent(events.CreateProcessUnknownEvent(progName, groupName, p.state.String()))
	}

	p.state = procState
}
func (p *Process) getExitCode() (int, error) {
	if p.cmd.ProcessState == nil {
		return -1, fmt.Errorf("no exit code")
	}
	if status, ok := p.cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		return status.ExitStatus(), nil
	}

	return -1, fmt.Errorf("no exit code")

}

func (p *Process) Stop(wait bool) {
	p.lock.Lock()
	p.stopByUser = true
	p.lock.Unlock()
	log.WithFields(log.Fields{"programM": p.GetName()}).Warn("stop the programM")
	sigs := strings.Fields("")
	waitsecs := 10 * time.Second
	stopasgroup := false
	killasgroup := stopasgroup
	if stopasgroup && !killasgroup {
		log.WithFields(log.Fields{"programM": p.GetName()}).Error("Cannot set stopasgroup=true and killasgroup=false")
	}

	go func() {
		stopped := false
		for i := 0; i < len(sigs) && !stopped; i++ {
			// send signal to process
			sig, err := signals.ToSignal(sigs[i])
			if err != nil {
				continue
			}
			log.WithFields(log.Fields{"programM": p.GetName(), "signal": sigs[i]}).Info("send stop signal to programM")
			p.Signal(sig, stopasgroup)
			endTime := time.Now().Add(waitsecs)
			//wait at most "stopwaitsecs" seconds for one signal
			for endTime.After(time.Now()) {
				//if it already exits
				if p.state != STARTING && p.state != RUNNING && p.state != STOPPING {
					stopped = true
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
		if !stopped {
			log.WithFields(log.Fields{"programM": p.GetName()}).Warn("force to kill the programM")
			p.Signal(syscall.SIGKILL, killasgroup)
		}
	}()
	if wait {
		for {
			// if the programM exits
			p.lock.RLock()
			if p.state != STARTING && p.state != RUNNING && p.state != STOPPING {
				p.lock.RUnlock()
				break
			}
			p.lock.RUnlock()
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *Process) GetStatus() string {
	if p.cmd.ProcessState.Exited() {
		return p.cmd.ProcessState.String()
	}
	return "running"
}
func (p *Process) Signal(sig os.Signal, sigChildren bool) error {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.sendSignal(sig, sigChildren)
}
func (p *Process) sendSignal(sig os.Signal, sigChildren bool) error {
	if p.cmd != nil && p.cmd.Process != nil {
		err := signals.Kill(p.cmd.Process, sig, sigChildren)
		return err
	}
	return fmt.Errorf("process is not started")
}
