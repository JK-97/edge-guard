package config

import "jxcore/lowapi/utils"

// 恢复的目录：默认配置备份路径 -> 配置路径
var mapSettings = map[string]string{
	"/edge/recover/jxcore_settings":  "/edge/jxcore/bin/settings",
	"/edge/recover/consul_conf.json": "/data/edgex/consul/config/consul_conf.json",
	"/edge/recover/interfaces":       "/etc/network/interfaces",
}

// 恢复系统默认配置
func ResetSystemConfig() error {
	for src, dst := range mapSettings {
		err := utils.CopyTo(src, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO 制定备份文件方案
func EnsureSystemConfigBackup() error {
	for src, dst := range mapSettings {
		if !utils.FileExists(src) {
			err := utils.CopyTo(dst, src)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
