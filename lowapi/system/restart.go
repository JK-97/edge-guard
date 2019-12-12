package system

import (
	"jxcore/gateway/log"
	"syscall"
	"time"
)

func RebootAfter(t time.Duration) error {
	log.Info("System will reboot in ", t)
	time.Sleep(t)
	return RunCommand("reboot")
}

func RestartJxcoreAfter(t time.Duration) {
	log.Info("Jxcore will restart in ", t)
	time.Sleep(t)
	syscall.Exit(0)
}
