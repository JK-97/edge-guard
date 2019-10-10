package signal

import (
	"os"
	"syscall"
)

func init() {

	// HUP     1    终端断线
	// INT     2    中断（同 Ctrl + C）
	// QUIT    3    退出（同 Ctrl + \）
	// TERM   15    终止
	// CONT   18    继续（与STOP相反， fg/bg命令）
	// STOP   19    暂停（同 Ctrl + Z）
	// KILL    9    强制终止
	StopSignals = []os.Signal{
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
		syscall.Signal(0x40),
	}

	HangUpSignals = []os.Signal{
		syscall.SIGHUP,
	}
}
