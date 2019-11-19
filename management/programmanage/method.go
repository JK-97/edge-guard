package programmanage

import (
	log "gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
	"strings"
)

func GetJxConfig() (config string) {
	return ProgramMconfig + ProgramSetting + BaseDepend
}

func GetBaseConfig() (config string) {
	return ProgramMconfig
}

// GetMucConfig 获取 mcu 相关程序的配置
func GetMcuConfig() string {
	return Watchdog + Mcuserver
}

// GetJxserving 获取 jxserving 配置
func GetJxserving() string {
	return Jxserving
}

func AddToStart(programName string) {
	if _, ok := ProgramCfgMap[programName]; ok == true {
		ProgramSetting = ProgramSetting + ProgramCfgMap[programName]
	} else {
		log.Info("this jxcore version does not suppoted this commponent,please update")
	}

}

func AddDependStart(programName string) {
	if _, ok := ProgramCfgMap[programName]; ok == true {
		lines := strings.Split(ProgramCfgMap[programName], "\n")
		newlines := make([]string, 0)
		newlines = append(newlines, lines[0], DependOnBase)
		for _, str := range lines[1:] {
			newlines = append(newlines, str)
		}

		ProgramCfgDepended := strings.Join(newlines, "\n")

		ProgramSetting = ProgramSetting + ProgramCfgDepended
	} else {
		log.Info("Missibg " + programName + " config,maybe jxcore low version")
	}

}
