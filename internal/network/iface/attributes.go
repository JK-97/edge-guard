package iface

import (
	"jxcore/config/yaml"
	"time"
)

const (
	testIP       = "114.114.114.114"
	highPriority = 5
)

var (
	currentIFace           string
	checkBestIFaceInterval = time.Second * 5
	ifacePriority          = yaml.Config.IFace.Priority
	backupIFace            = yaml.Config.IFace.Backup
)

func init() {
	interval, err := time.ParseDuration(yaml.Config.IFace.SwitchInterval)
	if err == nil {
		checkBestIFaceInterval = interval
	}
}
