package docker

import (
	"context"
	"fmt"
	"jxcore/journal"
	"jxcore/log"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func init() {
	journal.RegisterWorker(journal.ModeContainer, new(LoggerWorker))
}

// Collector 日志采集器
type Collector struct {
	Client    *client.Client
	Container *types.Container
}

// Collect 采集指定区间的日志
func (c *Collector) Collect(since, until int64) journal.BufferedLoggerWrapper {
	ctx := context.Background()
	st := time.Unix(since, 0).UTC()
	reader, err := c.Client.ContainerLogs(ctx, c.Container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      st.Format(time.RFC3339),
		Timestamps: true,
		Follow:     false,
		// Details:    true,
	})

	if err != nil {
		return nil
	}
	defer reader.Close()
	name := fmt.Sprintf("%s-%d-%d.log", strings.Trim(c.Container.Names[0], "/"), since, until)
	wrapper := journal.NewStreamLoggerWrapper(reader, name)
	return wrapper
}

// LoggerWorker 日志打包类
type LoggerWorker struct {
	*journal.TimeRange
	Client  *client.Client
	Include []string
	Exclude []string
}

// Mode 返回支持的日志模式
func (a *LoggerWorker) Mode() string {
	return journal.ModeContainer
}

// FetchAsync 异步采集日志
func (a *LoggerWorker) FetchAsync() (c chan journal.BufferedLoggerWrapper, wg *sync.WaitGroup) {
	f := filters.NewArgs()
	for _, inc := range a.Include {
		f.Add("name", inc)
	}
	containers, err := a.Client.ContainerList(context.Background(), types.ContainerListOptions{
		Quiet:   true,
		Filters: f,
	})
	if err != nil {
		log.Warn(err)
		return
	}

	c = make(chan journal.BufferedLoggerWrapper, 4)
	wg = new(sync.WaitGroup)
	for _, container := range containers {

		containerCopy := container
		collector := Collector{
			Client:    a.Client,
			Container: &containerCopy,
		}
		wg.Add(1)
		go func(c chan<- journal.BufferedLoggerWrapper, col *Collector) {
			c <- col.Collect(a.Since, a.Until)
		}(c, &collector)
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

	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	a.Client = cli
	d := config["docker-include"]
	if d != nil {
		a.Include = d.([]string)
	}
	d = config["docker-exclude"]
	if d != nil {
		a.Exclude = d.([]string)
	}
	return nil
}
