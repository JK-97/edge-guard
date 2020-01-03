package system

import (
	"fmt"
	"jxcore/lowapi/logger"
	"os/exec"

	"github.com/pkg/errors"
)

func RunCommand(command string) error {
	logger.Info("run command: ", command)
	output, err := exec.Command("/bin/bash", "-c", command).CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "output: "+string(output))
	}
	return err
}

func StopDisableService(name string) error {
	err1 := RunCommand("systemctl stop " + name)
	err2 := RunCommand("systemctl disable " + name)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("stop error: %v, disable error: %v", err1, err2)
	}
	return nil
}

func StartEnableService(name string) error {
	err1 := RunCommand("systemctl start " + name)
	err2 := RunCommand("systemctl enable " + name)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("start error: %v, enable error: %v", err1, err2)
	}
	return nil
}
