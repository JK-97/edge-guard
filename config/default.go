package config

import (
	"github.com/spf13/viper"
)

func SetDefault() {
	// viper.SetDefault("jxserving", true)
	// viper.SetDefault("fsindex", true)
	// viper.SetDefault("telegraf", true)
	// viper.SetDefault("powermanagement", true)
	// viper.SetDefault("mcuserver", true)
	// viper.SetDefault("db", true)
	// viper.SetDefault("tsdb", true)
	// viper.SetDefault("mq", true)
	// viper.SetDefault("fs", true)

	viper.SetDefault("components", map[string]bool{
		"jxserving":        true,
		"fsindex":          true,
		"telegraf":         true,
		"powermanagerment": true,
		"mcuserver":        true,
		"db":               true,
		"tsdb":             true,
		"mq":               true,
		"fs":               true,
	})
	viper.SetDefault("fixedresolver", "")
	viper.SetDefault("debug", true)

	viper.SetDefault("priority", []string{
		"eth0",
		"usb0",
		"eth1",
		"usb1"})
	viper.SetDefault("backup", "usb0")
	viper.SetDefault("switch_interval", "5s")

	viper.SetDefault("mount_cfg", map[string]string{
		"/dev/mmcblk1p1": "/media/card",
	})
	viper.SetDefault("port", ":80")
}
