package logs

import (
	"fmt"
	"jxcore/oplog/types"
	"strings"
	"time"
)

const indent = "   "

type LogMessage struct {
	recordTime  time.Time
	messageType string
	description string
	types.Oplog
}

func NewOplog(messageType, description string) types.Oplog {
	return &LogMessage{
		recordTime:  time.Now(),
		messageType: messageType,
		description: description,
	}
}

func (l *LogMessage) GetRecordTime() time.Time {
	return l.recordTime
}

func (l *LogMessage) GetMessageType() string {
	return l.messageType
}

func (l *LogMessage) GetDescription() string {
	return l.description
}

func (l *LogMessage) Marshal() []byte {
	return []byte(fmt.Sprintf(strings.Join([]string{"%s", "%s", "%s\n"}, indent), l.recordTime.Format("Mon Jan 2 15:04:05 MST 2006"), l.messageType, l.description))
}

func (l *LogMessage) UnMarshal(data []byte) error {
	info := strings.Split(string(data), indent)
	time, err := time.Parse("Mon Jan 2 15:04:05 MST 2006", info[0])
	if err != nil {
		return err
	}
	l.recordTime = time
	l.messageType = info[1]
	l.description = info[2]
	return nil
}
