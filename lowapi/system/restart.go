package system

import (
	"syscall"
	"time"

	"github.com/JK-97/edge-guard/lowapi/logger"
)

func RebootAfter(t time.Duration) error {
	logger.Info("System will reboot in ", t)
	time.Sleep(t)
	return RunCommand("reboot")
}

func RestartEdgeguardAfter(t time.Duration) {
	logger.Info("edge-guard will restart in ", t)
	time.Sleep(t)
	syscall.Exit(0)
}
