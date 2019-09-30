// +build cgo

package systemd

import (
	"fmt"
	"jxcore/journal"
	"jxcore/log"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
)

func init() {
	journal.RegisterWorker(journal.ModeKernel, new(LoggerWorker))
}

// Collector 日志采集器
type Collector struct {
	// Journal *sdjournal.Journal
	Service string
}

// Collect 采集指定区间的日志
func (c *Collector) Collect(since, until int64) journal.BufferedLoggerWrapper {
	matches := []sdjournal.Match{}
	if c.Service != "" {
		m := sdjournal.Match{
			Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
			Value: c.Service,
		}
		matches = append(matches, m)
	}

	config := sdjournal.JournalReaderConfig{
		Since:   -time.Now().Sub(time.Unix(since, 0)),
		Matches: matches,
	}

	r, err := sdjournal.NewJournalReader(config)

	if err != nil {
		log.Warn("Error opening journal: %s", err)
		return nil
	}

	if r == nil {
		log.Warn("Got a nil reader")
		return nil
	}

	defer r.Close()
	var name string
	if c.Service == "" {
		name = fmt.Sprintf("%s-%d-%d.log", "kernel", since, until)

	} else {
		name = fmt.Sprintf("%s-%d-%d.log", c.Service, since, until)
	}
	return journal.NewStreamLoggerWrapper(r, name)
}

// LoggerWorker 日志打包类
type LoggerWorker struct {
	*journal.TimeRange

	Include []string
}

// Mode 返回支持的日志模式
func (a *LoggerWorker) Mode() string {
	return journal.ModeKernel
}

// FetchAsync 异步采集日志
func (a *LoggerWorker) FetchAsync() (c chan journal.BufferedLoggerWrapper, wg *sync.WaitGroup) {

	c = make(chan journal.BufferedLoggerWrapper, 4)
	wg = new(sync.WaitGroup)
	if len(a.Include) > 0 {
		for _, srv := range a.Include {
			if !strings.HasSuffix(srv, ".service") {
				srv = srv + ".service"
			}
			collector := Collector{Service: srv}
			wg.Add(1)
			go func(c chan<- journal.BufferedLoggerWrapper, col *Collector) {
				c <- col.Collect(a.Since, a.Until)
			}(c, &collector)
		}
	} else {
		wg.Add(1)
		collector := Collector{}
		c <- collector.Collect(a.Since, a.Until)
	}

	return
}

// CanAppend 指示是否支持追加写入
func (a *LoggerWorker) CanAppend() bool {
	return true
}

// InitConfig 初始化配置
func (a *LoggerWorker) InitConfig(timeRange *journal.TimeRange, config map[string]interface{}) error {
	a.TimeRange = timeRange

	include := config["kernel-include"]
	if include != nil {
		a.Include = include.([]string)
	}

	return nil
}
