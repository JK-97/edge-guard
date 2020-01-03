package iface

import (
	// "jxcore/config/yaml"
	"github.com/spf13/viper"
	"time"
)

const (
	testIP       = "114.114.114.114"
	highPriority = 5
)

var (
	currentIFace           string
	checkBestIFaceInterval = time.Second * 5
	ifacePriority          = viper.GetStringMapString("backup")
	backupIFace            = viper.GetString("backup")
)

func init() {
	interval, err := time.ParseDuration(viper.GetString("switch_interval"))
	if err == nil {
		checkBestIFaceInterval = interval
	}
}
