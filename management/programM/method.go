package programM

import "jxcore/log"

func GetJxConfig() (config string) {
    return ProgramMconfig+ BaseDepend +ProgramSetting
}

func GetBaseConfig() (config string) {
    return ProgramMconfig
}

func AddToStart(programName string) {
    if _, ok := ProgramCfgMap[programName]; ok == true {
        ProgramSetting = ProgramSetting + ProgramCfgMap[programName]   
    }else{
        log.Info("this jxcore  version does not suppoted this commponent,please update")
    }
    
}


func AddDependStart(programName string){
    if _, ok := ProgramCfgMap[programName]; ok == true {
        ProgramSetting = ProgramSetting + ProgramCfgMap[programName]+DependOnBase
    }else{
        log.Info("this jxcore  version does not suppoted this commponent,please update")
    }
    
    
}