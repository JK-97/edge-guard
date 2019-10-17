package core

import "sync"

var DnsOnce sync.Once
var jxcore *JxCore
var lock *sync.Mutex = &sync.Mutex{}


const (
    
    logBase  = "/edge/logs/"
    YamlComponentSetting = "/edge/jxcore/bin/settings.yaml"
)

