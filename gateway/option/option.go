package option

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

// Service 服务配置
type Service struct {
	// LocalOnly 服务是否只存在节点上
	LocalOnly bool
	// MasterOnly 服务是否只存在云端
	MasterOnly bool
	// 转发 URL
	Proxy string
}

// Route 服务路由
type Route struct {
	// 服务名
	Name string
	// URL 字符串替换模式
	ReplaceMent string `toml:"Replace"`

	// 正则表达式
	Matcher string
}

// ListenConfig 服务监听地址的配置
type ListenConfig struct {
	// Addr 监听地址
	Addr string
	// ExtraAddrs 额外的监听地址
	ExtraAddrs []string
	// SocketAddr 监听 Unix Socket 地址
	SocketAddr string
	// SocketMode 文件模式，八进制 0777
	SocketMode uint
}

// MessageQueueConfig 消息队列相关的配置
type MessageQueueConfig struct {
	// MessageQueueURI 本地 MQ URI
	MessageQueueURI    string `toml:"MessageQueue"`
	MessageQueuePrefix string
}

// ProxyServerConfig 代理服务配置
type ProxyServerConfig struct {
	MasterProxy string

	Services map[string]Service
	Routes   []Route
}

// ContainerConfig 容器相关配置
type ContainerConfig struct {

	// Docker 域 unix:///var/run/docker.sock 或 http://127.0.0.1:2375
	DockerDomin string
	// ComposeBinary docker-compose 可执行文件路径
	ComposeBinary string
	// Docker Compose 使用的父目录路径
	ComposeBaseDir string
}

// StorageConfig 存储配置
type StorageConfig struct {
	TempDir   string
	LocalDir  string
	RemoteDir string
	TempFSDir string

	CephScript string
	MasterIP   string

	BatchSize       int
	RetryInterval   int
	RecheckInterval int
}

// ConfigAgentConfig ConfigAgent 相关的配置
type ConfigAgentConfig struct {
	Timeout time.Duration
	Host    string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	ListenConfig

	MessageQueueConfig

	// WorkingDirectory 工作目录
	WorkingDirectory string `toml:"Dir"`
	// AiServingURL AI Serving 地址
	AiServingURL string `toml:"AiServing"`
	// EgdeX 代理地址
	EgdeXURL string `toml:"JXEdge"`

	// 是否允许服务动态注册
	EnableDynamicService bool

	// ContainerConfig 容器服务相关配置
	ContainerConfig

	// Proxy 代理转发服务配置
	ProxyServerConfig

	ConfigAgent ConfigAgentConfig

	Storage StorageConfig

	Device DeviceConfig

	Database DatabaseConfig

	TimeSeries TimeSeriesDBConfig
}

// TimeSeriesDBConfig 时序数据库相关配置
type TimeSeriesDBConfig struct {
	InfluxDBPort     int // InfluxDB 使用的端口号
	InfluxDBUser     string
	InfluxDBPassword string
	StatdPort        int // Statsd 使用的端口号
	StorePrefix      string
}

// Decode 从配置文件字符串中加载服务配置
func Decode(t string, c *ServerConfig) (toml.MetaData, error) {
	meta, err := toml.Decode(t, c)
	if err != nil {
		return meta, err
	}

	if c.SocketMode > 0 {
		m := c.SocketMode
		if m < 0700 {
			// 权限不得小于 0700
			c.SocketMode = 0700
		} else if m >= 700 && m < 1000 {
			// 八进制数被当作十进制数解析
			var a uint
			a = ((m / 100) * 0100) + ((m / 10) % 10 * 010) + (m % 10)
			c.SocketMode = a
		}
	}

	return meta, err
}

// ServerConfigFromFile 从指定的配置文件路径中加载服务配置
func ServerConfigFromFile(f string, c *ServerConfig) (toml.MetaData, error) {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return toml.MetaData{}, err
	}
	return Decode(string(bs), c)
}

// SaveToWriter 序列化至指定的 Writer
func (c *ServerConfig) SaveToWriter(w io.Writer) error {
	encoder := toml.NewEncoder(w)

	return encoder.Encode(c)
}

// DefaultListenConfig 默认的监听配置
func DefaultListenConfig() ListenConfig {
	return ListenConfig{
		Addr: ":8080",
	}
}

// DefaultContainerConfig 默认的容器配置
func DefaultContainerConfig() ContainerConfig {
	return ContainerConfig{
		DockerDomin:    "unix:///var/run/docker.sock",
		ComposeBinary:  "docker-compose",
		ComposeBaseDir: "/tmp/compose",
	}
}

// DefaultServerConfig 默认的服务配置
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		ListenConfig:    DefaultListenConfig(),
		ContainerConfig: DefaultContainerConfig(),
		ProxyServerConfig: ProxyServerConfig{
			MasterProxy: "http://10.201.0.1:9000",
			Services:    make(map[string]Service),
			Routes:      make([]Route, 0),
		},
		TimeSeries: TimeSeriesDBConfig{
			InfluxDBPort:     8086,
			InfluxDBUser:     "root",
			InfluxDBPassword: "root",
			StatdPort:        8125,
		},
	}
}

// DeviceConfig 设备相关的配置
type DeviceConfig struct {
	Timeout time.Duration
	Hosts   map[string]string
}

// DefaultDeviceConfig 默认的设备相关的配置
func DefaultDeviceConfig() DeviceConfig {
	return DeviceConfig{
		Timeout: 1 * time.Second,
		Hosts: map[string]string{
			"core-command": "http://127.0.0.1:48082",
		},
	}
}

// DatabaseConfig 数据库相关的配置
type DatabaseConfig struct {
	Timeout time.Duration
	Host    string
}

// DefaultDatabaseConfig 默认的数据库相关的配置
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Timeout: 1 * time.Second,
		Host:    "http://127.0.0.1:9998",
	}
}
