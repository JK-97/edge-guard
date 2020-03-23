package oplog

import (
	"jxcore/oplog/types"
	"time"
)

func DefaultTimeFilter(from, until time.Time) types.FilterFunc {
	return func(o types.Oplog) bool {
		recordTime := o.GetRecordTime()
		if !(recordTime.Before(until) && recordTime.After(from)) {
			return false
		}
		return true
	}
}

func DefaultTypeFilter(logMessageType string) types.FilterFunc {
	return func(o types.Oplog) bool {
		if logMessageType == "all" {
			return true
		}
		if !(o.GetMessageType() == logMessageType) {
			return false
		}
		return true
	}
}
