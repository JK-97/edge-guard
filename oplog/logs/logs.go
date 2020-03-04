package logs

import (
	"encoding/json"
	"fmt"
	"jxcore/oplog/types"
	"strings"
	"time"
)

const indent = "   "

type LogMessage struct {
	recordTime  time.Time `json:"recordtime"`
	messageType string    `json:"messagetype"`
	description string    `json:"description"`
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
	return []byte(fmt.Sprintf(strings.Join([]string{"%s", "%s", "%s\n"}, indent), l.recordTime.Format("2006-01-02 15:04:05"), l.messageType, l.description))
}

func (l *LogMessage) UnMarshal(data []byte) error {
	info := strings.Split(string(data), indent)
	time, err := time.Parse("2006-01-02 15:04:05", info[0])
	if err != nil {
		return err
	}
	l.recordTime = time
	l.messageType = info[1]
	l.description = info[2]
	return nil
}

func (l *LogMessage) String() string {
	return string(l.Marshal())
}
func (l *LogMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		RecordTime  string `json:"recordtime"`
		MessageType string `json:"messagetype"`
		Description string `json:"description"`
	}{
		RecordTime:  l.GetRecordTime().Format("2006-01-02 15:04:05"),
		MessageType: l.GetMessageType(),
		Description: l.GetDescription(),
	})
}
