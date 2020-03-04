package manager

import (
	"jxcore/oplog/types"
)

func init() {
	This = NewManager(defaultLogPath)
}

var defaultLogPath = "/var/log/jxcore_event.log"
var This *Manager = nil

func Insert(o types.Oplog) error {
	return This.Insert(o)
}
func Remove(o types.Oplog) error {
	return This.Remove(o)
}

func ListAll() ([]types.Oplog, error) {
	return This.ListAll()
}

func FindMany(f ...types.FilterFunc) ([]types.Oplog, error) {
	return This.FindMany(f...)
}
func GetLogFileName() string {
	return This.logFile.Name()
}
