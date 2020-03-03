package oplog

import (
	"jxcore/oplog/types"
	"jxcore/oplog/manager"
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
