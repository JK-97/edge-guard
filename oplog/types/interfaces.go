package types

import (
	"time"
)

type Oplog interface {
	GetRecordTime() time.Time
	GetMessageType() string
	GetDescription() string
	Marshal() []byte
	UnMarshal([]byte) error
}

type FilterFunc func(Oplog) bool

type filter interface {
	Filter(Oplog) bool
}

func (f FilterFunc) Filter(oplog Oplog) bool {
	return f(oplog)
}

type OplogManager interface {
	Insert(Oplog) error
	Remove(Oplog) error
	ListAll() ([]Oplog, error)
	FindMany(...FilterFunc) ([]Oplog, error)
}
