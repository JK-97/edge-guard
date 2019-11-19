package core

import "sync"

var DnsOnce sync.Once
var jxcore = NewJxCore()

const (
	logBase = "/edge/logs/"
)
