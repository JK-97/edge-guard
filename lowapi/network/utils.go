package network

import (
	"os/exec"
)

func LockResolvConf() error {
	return exec.Command("chattr", "+i", "/etc/resolv.conf").Run()
}

func UnlockResolvConf() error {
	return exec.Command("chattr", "-i", "/etc/resolv.conf").Run()
}
