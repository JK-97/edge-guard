package system

import (
	log "jxcore/lowapi/logger"
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
