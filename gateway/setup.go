package gateway

import (
	"jxcore/gateway/option"
	"jxcore/gateway/store"
	"jxcore/lowapi/logger"
	"os"
	"path/filepath"
)

var (
	workingDir string = "/data/local/gateway"
	configPath string = "/edge/jxcore/bin/gateway.cfg"
)

// Setup 配置 Gateway
func Setup() {
	if _, err := option.ServerConfigFromFile(configPath, &ServerOptions); err != nil {
		if os.IsNotExist(err) {
			ServerOptions = option.DefaultServerConfig()
			logger.Infof("File Not Found Use Default Config\n")
			ServerOptions.SaveToWriter(os.Stdout)
		} else {
			logger.Error(err)
		}
	}

	if fileInfo, err := os.Stat(workingDir); os.IsNotExist(err) {
		// WorkingDir 不存在，尝试创建
		err = os.MkdirAll(workingDir, os.ModePerm)
		if err != nil {
			logger.Infof("Error %s\n", err)
		}

	} else if !fileInfo.IsDir() {
		logger.Errorf("[%v] is not directory\n", workingDir)
	} else {
		defaultStore = store.NewLevelDBStore(filepath.Join(workingDir, "gateway.db"))
		logger.Info("Working directory: ", workingDir)
	}
}
