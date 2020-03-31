package oplog

import (
	"github.com/JK-97/edge-guard/oplog/manager"
	"github.com/JK-97/edge-guard/oplog/types"
)

func Insert(o types.Oplog) error {
	return manager.Insert(o)
}

func Remove(o types.Oplog) error {
	return manager.Remove(o)
}

func ListAll() ([]types.Oplog, error) {
	return manager.ListAll()
}

func FindMany(f ...types.FilterFunc) ([]types.Oplog, error) {
	return manager.FindMany(f...)
}

func Marshal(o types.Oplog) []byte {
	return o.Marshal()
}

func UnMarshal(data []byte, o types.Oplog) error {
	return o.UnMarshal(data)
}

func GetLogFileName() string {
	return manager.GetLogFileName()
}
