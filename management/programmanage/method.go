package programmanage

import (
	"strings"

	log "github.com/JK-97/edge-guard/lowapi/logger"
)

func GetJxConfig() (config string) {
	return ProgramMconfig + ProgramSetting + BaseDepend
}

func GetBaseConfig() (config string) {
	return ProgramMconfig
}

func AddToStart(programName string) {
	if _, ok := ProgramCfgMap[programName]; ok {
		ProgramSetting = ProgramSetting + ProgramCfgMap[programName]
	} else {
		log.Info("this edge-guard version does not suppoted this commponent,please update")
	}

}

func AddDependStart(programName string) {
	if _, ok := ProgramCfgMap[programName]; ok {
		lines := strings.Split(ProgramCfgMap[programName], "\n")
		newlines := make([]string, 0)
		newlines = append(newlines, lines[0], DependOnBase)
		for _, str := range lines[1:] {
			newlines = append(newlines, str)
		}

		ProgramCfgDepended := strings.Join(newlines, "\n")

		ProgramSetting = ProgramSetting + ProgramCfgDepended
	} else {
		log.Info("Missibg " + programName + " config,maybe edge-guard low version")
	}

}
