package journal

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

// 日志模式
const (
	ModeContainer = "container" // 容器
	ModeKernel    = "kernel"    // 内核
	ModeService   = "service"   // 服务
)

// MetaFileItem 元数据信息
type MetaFileItem struct {
	TimeRange
	// StartTime int64
	// EndTime   int64
}

// MetaFileItemSlice MetaFileItem 切片
type MetaFileItemSlice []MetaFileItem

func (a MetaFileItemSlice) Len() int      { return len(a) }
func (a MetaFileItemSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a MetaFileItemSlice) Less(i, j int) bool {
	if a[i].Since == a[j].Since {
		return a[i].Until < a[j].Until
	}
	return a[i].Since < a[j].Since
}

// LogArchiveMeta 日志文件元数据
type LogArchiveMeta struct {
	Files      []string                `json:"files"`
	Rotates    map[string]MetaFileItem `json:"rotates"`
	CreateTime int64                   `json:"create_time"`
	ModifyTime int64                   `json:"modify_time"`
	Mode       string                  `json:"mode"`
	Day        int                     `json:"day"` // 元数据所属的日志 YYYYMMDD 格式
}

// Deserialize 从 Reader 中读取
func (m *LogArchiveMeta) Deserialize(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, m)
	return err
}

// Serialize 写入指定的 writer
func (m *LogArchiveMeta) Serialize(w io.Writer) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// LoggerWrapper 将日志信息保存在 io.Writer 上
type LoggerWrapper interface {
	// Save 将日志信息保存在 io.Writer 上
	Save(w io.Writer) error
}

// BufferedLoggerWrapper 将日志信息保存在 io.Writer 上
type BufferedLoggerWrapper interface {
	LoggerWrapper
	// Bytes 获取可写入的数据
	Bytes() []byte
	// Name 文件名
	Name() string
}

// TimeRange 时间范围
type TimeRange struct {
	Since int64
	Until int64
}

// LoggerCollector 日志采集器
type LoggerCollector interface {
	// Collect 采集指定区间的日志
	Collect(since, until int64) BufferedLoggerWrapper
}

// LoggerWorker 日志收集器
type LoggerWorker interface {
	// FetchAsync 异步采集日志
	FetchAsync() (c chan BufferedLoggerWrapper, wg *sync.WaitGroup)
	// InitConfig 初始化配置
	InitConfig(timeRange *TimeRange, config map[string]interface{}) error

	// CanAppend 指示是否支持追加写入
	CanAppend() bool
}

// LoggerCleanWorker 有垃圾需要清理
type LoggerCleanWorker interface {
	LoggerWorker

	// Clean 文件清理
	Clean(ttl time.Duration) error
}

// LogArchiveWorker 日志打包类
type LogArchiveWorker interface {
	// Mode 返回支持的日志模式
	Mode() string

	// MetaFile 获取元数据
	MetaFile() LogArchiveMeta

	// SetMetaFile 设置元数据
	SetMetaFile(LogArchiveMeta)

	// OpenArchive 读取文件
	OpenArchive(r io.Reader) error

	// MakeArchive 打包文件
	MakeArchive(w io.Writer) error

	// // AppendFiles 添加日志文件
	// AppendFiles(v ...interface{}) error
}

// StreamLoggerWrapper 可重复读取的 LoggerWrapper 实现
type StreamLoggerWrapper struct {
	// Reader io.Reader
	reader *bytes.Buffer
	name   string
}

// NewStreamLoggerWrapper 由 bytes.Buffer 实现的 StreamLoggerWrapper
func NewStreamLoggerWrapper(r io.Reader, name string) *StreamLoggerWrapper {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}
	return &StreamLoggerWrapper{
		reader: bytes.NewBuffer(buf),
		name:   name,
	}
}

// Save 将日志信息保存在 io.Writer 上
func (log *StreamLoggerWrapper) Save(w io.Writer) error {
	_, err := w.Write(log.reader.Bytes())
	return err
}

// Bytes 获取可写入的数据
func (log *StreamLoggerWrapper) Bytes() []byte {
	if log == nil {
		return nil
	}
	return log.reader.Bytes()
}

// Name 文件名
func (log *StreamLoggerWrapper) Name() string {
	if log == nil {
		return ""
	}
	return log.name
}

// RegisteredWorkers 已注册的 Worker
var RegisteredWorkers map[string]LoggerWorker = make(map[string]LoggerWorker, 0)

// RegisterWorker 注册 Worker
func RegisterWorker(mode string, worker LoggerWorker) {
	RegisteredWorkers[mode] = worker
}

// RegisteredModes 获取已注册的插件
func RegisteredModes() []string {
	result := make([]string, 0, 3)

	for key := range RegisteredWorkers {
		result = append(result, key)
	}

	return result
}
