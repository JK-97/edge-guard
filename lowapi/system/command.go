package system

import (
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
