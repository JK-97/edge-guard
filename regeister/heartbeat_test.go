package regeister

import (
	"fmt"
	"jxcore/log"
	"os/exec"
	"testing"
)

func TestGetMyIP(t *testing.T) {
	out, err := exec.Command("/bin/bash", "-c", "rm -r /home/marshen/mongo/*").Output()
	if err != nil {
		log.Error(err)
	}
	fmt.Println(out)
}
