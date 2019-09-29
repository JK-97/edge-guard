package rfile

import (
	"io/ioutil"
	"jxcore/journal"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func init() {
	journal.RegisterWorker(journal.ModeService, new(LoggerWorker))
}

const suffixFormat string = "2006-01-02.log"

// GetSuffix 获取后缀
func GetSuffix(t int64) string {
	return time.Unix(t, 0).Local().Format(suffixFormat)
}

// Collector 日志采集器
type Collector struct {
	FileName string
}

// Collect 采集指定区间的日志
func (c *Collector) Collect(since, until int64) journal.BufferedLoggerWrapper {
	reader, err := os.Open(c.FileName)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer reader.Close()
	name := strings.Trim(filepath.Base(c.FileName), "/")
	wrapper := journal.NewStreamLoggerWrapper(reader, name)
	return wrapper
}

// LoggerWorker 日志打包类
type LoggerWorker struct {
	*journal.TimeRange
	Directorys []string
}

// Mode 返回支持的日志模式
func (a *LoggerWorker) Mode() string {
	return journal.ModeService
}

// FetchAsync 异步采集日志
func (a *LoggerWorker) FetchAsync() (c chan journal.BufferedLoggerWrapper, wg *sync.WaitGroup) {

	c = make(chan journal.BufferedLoggerWrapper, 4)
	wg = new(sync.WaitGroup)

	suffix := GetSuffix(a.Since)
	for _, directory := range a.Directorys {
		err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, suffix) {
				wg.Add(1)
				collector := Collector{
					FileName: path,
				}
				go func(c chan<- journal.BufferedLoggerWrapper, col *Collector) {
					c <- col.Collect(a.Since, a.Until)
				}(c, &collector)
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}

	return
}

// CanAppend 指示是否支持追加写入
func (a *LoggerWorker) CanAppend() bool {
	return false
}

// Clean 文件清理
func (a *LoggerWorker) Clean(ttl time.Duration) error {
	deadLine := time.Now().Add(-ttl).Unix()

	for _, folder := range a.Directorys {
		if infos, err := ioutil.ReadDir(folder); err == nil {
			for _, info := range infos {
				mod := info.ModTime()
				if mod.Unix() < deadLine {
					os.Remove(info.Name())
				}
			}
		}
	}
	return nil
}

// InitConfig 初始化配置
func (a *LoggerWorker) InitConfig(timeRange *journal.TimeRange, config map[string]interface{}) error {
	a.TimeRange = timeRange

	d := config["rotate-directory"]
	if d != nil {
		a.Directorys = d.([]string)
	}
	return nil
}
