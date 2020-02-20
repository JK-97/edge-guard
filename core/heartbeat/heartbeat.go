package heartbeat

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"syscall"
	"time"

	"jxcore/core/heartbeat/collector"
	"jxcore/core/heartbeat/message"
	"jxcore/lowapi/logger"

	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

// 默认配置
var (
	DefaultInterval     = time.Second * 5
	DefaultWriteTimeout = time.Second * 5
)

var protocolHead = []byte{0x6C, 0xC6}

// Status 心跳状态
type Status int

// 心跳状态
const (
	Disconnected Status = iota // 未连接
	Connected                  // 连接成功
	Beat                       // 心跳包发送成功
)

// Option option
type Option byte

// 序列化方式
const (
	OptionPlain    Option = 0
	OptionProtoBuf Option = 1 << 5

	OptionRequireMask Option = OptionProtoBuf - 1
	OptionWorkerID    Option = 0x01 // 设备 ID
	OptionCPU         Option = 0x02 // CPU 信息
	OptionMemory      Option = 0x03 // 内存信息
	OptionNetwork     Option = 0x04 // 网络信息
	OptionDisk        Option = 0x05 // 磁盘信息
	OptionFileSystem  Option = 0x06 // 文件系统信息
	OptionPlatform    Option = 0x07 // 平台信息
)

// PacketHandler 处理服务端请求
type PacketHandler = func() ([]byte, Option)

// HeartBeater 心跳发送器
type HeartBeater struct {
	RemoteAddress   string                                // 远端地址
	AllowErrorCount int                                   // 允许出错数
	Interval        time.Duration                         // 心跳间隔
	WriteTimeout    time.Duration                         // 心跳发送超时时间
	Status          Status                                // 更新心跳状态
	Handlers        map[message.RequireType]PacketHandler // 处理从服务端接收的包
	InfoCollector   collector.InfoCollector               // 获取设备信息
	OnDNSError      func(err error) error                 // 再 DNS 解析出错时调用

	conn     net.Conn    // 心跳连接
	toBeSent chan []byte // 待发送队列
}

// NewHeartBeater HeartBeater for addr
func NewHeartBeater(addr string, collector collector.InfoCollector, AllowErrorCount int, OnDNSError func(err error) error) *HeartBeater {
	return &HeartBeater{
		RemoteAddress:   addr,
		AllowErrorCount: AllowErrorCount,
		Interval:        DefaultInterval,
		WriteTimeout:    DefaultWriteTimeout,
		Status:          Disconnected,
		Handlers:        make(map[message.RequireType]PacketHandler),
		InfoCollector:   collector,
		OnDNSError:      OnDNSError,

		conn:     nil,
		toBeSent: make(chan []byte, 6),
	}
}

// NewDefaultHeartBeater HeartBeater for addr with default handlers
func NewDefaultHeartBeater(addr string, AllowErrorCount int) *HeartBeater {
	beater := NewHeartBeater(addr, collector.NewLinuxCollector(), AllowErrorCount, nil)
	beater.RegisterDefaultHandlers()
	return beater
}

// Beat 向服务器发送心跳包
func (b *HeartBeater) Beat(ctx context.Context) error {
	for {
		errCount := 0
		// 连接 Device Manager
		logger.Info("Trying to connect device manager")
		for b.conn == nil {
			err := b.connect()
			if isDNSError(err) && b.OnDNSError != nil {
				err = b.OnDNSError(err)
			}
			if err != nil {
				errCount++
				if errCount > b.AllowErrorCount {
					return err
				}
				logger.Warn(err)
				time.Sleep(200 * time.Millisecond)
			}
		}

		logger.Info("Connected to device manager")
		if err := b.beat(ctx, b.conn); err != nil {
			logger.Warn(err)
		}

		if err := b.disconnect(); err != nil {
			logger.Error(err)
		}
	}
}

// connect 连接服务器
func (b *HeartBeater) connect() error {
	conn, err := net.DialTimeout("tcp", b.RemoteAddress, b.WriteTimeout)
	if err != nil {
		return err
	}

	b.conn = conn
	b.Status = Connected
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		setKeepaliveParameters(tcpConn)
	} else {
		logger.Error("not a tcp connection")
	}
	return err
}

func setKeepaliveParameters(conn *net.TCPConn) {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		logger.Warn("on getting raw connection object for keepalive parameter setting", err.Error())
	}

	err = rawConn.Control(
		func(fdPtr uintptr) {
			// got socket file descriptor. Setting parameters.
			fd := int(fdPtr)

			// 修复连接重连慢，原因是jxcore心跳包write无法发现tcp连接断开。
			// 设置ack timeout，使得如果write在10秒内没有收到ack，断开连接。
			// https://blog.cloudflare.com/when-tcp-sockets-refuse-to-die/
			timeoutSec := viper.GetInt("heartbeat_timeout_sec")
			// TCP_USER_TIMEOUT use millisecond
			err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, unix.TCP_USER_TIMEOUT, timeoutSec*1000)
			if err != nil {
				logger.Warn("on setting keepalive retry interval", err.Error())
			}
		})
	if err != nil {
		logger.Error(err)
	}
}

// disconnect 断开心跳连接
func (b *HeartBeater) disconnect() error {
	var err error
	if b.conn != nil {
		err = b.conn.Close()
		b.conn = nil
		b.Status = Disconnected
	}
	return err
}

// sendEmpty 发送空心跳包
func (b *HeartBeater) sendEmpty() {
	b.sendPacket(nil, OptionPlain)
}

// sendPacket 发送数据包
func (b *HeartBeater) sendPacket(p []byte, option Option) {
	packet := b.makePacket(p, byte(option))
	b.toBeSent <- packet
}

func (b *HeartBeater) beat(ctx context.Context, conn net.Conn) error {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := b.readFromConn(conn); err != nil {
			logger.Error(err)
		}
		cancel()
	}()

	// 第一个心跳包，先上报 WorkerID， 网络状态
	b.handleRequire(message.RequireType_RequireWorkerID, message.RequireType_RequireNetwork)

	ticker := time.NewTicker(b.Interval)
	defer ticker.Stop()

	errorCount := 0
	for {
		if errorCount >= b.AllowErrorCount {
			return ErrMaxErrorCountExceed
		}

		select {
		case p := <-b.toBeSent:
			_ = conn.SetWriteDeadline(time.Now().Add(b.WriteTimeout))
			_, err := conn.Write(p)
			if err != nil {
				if isConnectionBroken(err) {
					return err
				}

				logger.Error(err)
				errorCount++
				continue
			}

			b.Status = Beat
			// 发送成功，重置出错数量
			errorCount = 0
		case <-ticker.C:
			// 每隔固定时间间隔发送空数据包
			b.sendEmpty()
		case <-ctxWithCancel.Done():
			logger.Info("Heartbeat Canceled")
			return ctxWithCancel.Err()
		}
	}
}

func (b *HeartBeater) makePacket(p []byte, option byte) []byte {
	var buffer bytes.Buffer
	buffer.Write(protocolHead)
	if len(p) == 0 {
		// 心跳包，长度为 0
		buffer.WriteByte(0)
		buffer.WriteByte(0)
		return buffer.Bytes()
	}
	_p := make([]byte, 2)
	binary.LittleEndian.PutUint16(_p, uint16(len(p)))
	buffer.Write(_p)
	buffer.WriteByte(option)
	buffer.Write(p)

	return buffer.Bytes()
}
