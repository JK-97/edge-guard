package programmanage

import (
	"jxcore/log"
	"strings"
)

func GetJxConfig() (config string) {
	return ProgramMconfig + ProgramSetting // +BaseDepend
}

func GetBaseConfig() (config string) {
	return ProgramMconfig
}

func AddToStart(programName string) {
	if _, ok := ProgramCfgMap[programName]; ok == true {
		ProgramSetting = ProgramSetting + ProgramCfgMap[programName]
	} else {
		log.Info("this jxcore  version does not suppoted this commponent,please update")
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
		log.Info("this jxcore  version does not suppoted this commponent,please update")
	}

}
