package system

import (
	"jxcore/lowapi/logger"
	"syscall"
	"time"
)

func RebootAfter(t time.Duration) error {
	logger.Info("System will reboot in ", t)
	time.Sleep(t)
	return RunCommand("reboot")
}

func RestartJxcoreAfter(t time.Duration) {
	logger.Info("Jxcore will restart in ", t)
	time.Sleep(t)
	syscall.Exit(0)
}
