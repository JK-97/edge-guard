package oplog

import (
	"jxcore/oplog/types"
	"time"
)

func DefaultTimeFilter(o Oplog, from, until time.Time) types.FilterFunc {
	return func(oplog) bool {

	}
}
