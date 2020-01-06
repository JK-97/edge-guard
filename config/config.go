package config

import (
	"jxcore/lowapi/logger"

	"github.com/spf13/viper"
)

var CfgFile string

func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/")
		viper.SetConfigName("settings")
	}

	// viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			logger.Info(err)
		} else {
			logger.Info("Using config file:", viper.ConfigFileUsed())
		}
	}

}
