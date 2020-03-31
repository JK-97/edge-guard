package heartbeat

import (
	"bufio"
	"encoding/binary"
	"github.com/JK-97/edge-guard/core/heartbeat/message"
	"io"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/JK-97/go-utils/logger"
)

// RegisterDefaultHandlers 注册默认回调
func (b *HeartBeater) RegisterDefaultHandlers() {
	registerHandler := func(rt message.RequireType, pb proto.Message, optionType Option) {
		b.RegisterHandler(rt, func() (bytes []byte, option Option) {
			b, err := proto.Marshal(pb)
			if err != nil {
				logger.Error(err)
			}
			return b, OptionProtoBuf | optionType
		})
	}

	workerID := b.InfoCollector.GetWorkerID()
	registerHandler(message.RequireType_RequireWorkerID, &workerID, OptionWorkerID)

	cpuInfo := b.InfoCollector.GetCPUInfo()
	registerHandler(message.RequireType_RequireCPU, &cpuInfo, OptionCPU)

	memoryInfo := b.InfoCollector.GetMemoryInfo()
	registerHandler(message.RequireType_RequireMemory, &memoryInfo, OptionMemory)

	networks := b.InfoCollector.GetNetworkInfo()
	registerHandler(message.RequireType_RequireNetwork, &networks, OptionNetwork)

	diskInfo := b.InfoCollector.GetDiskStat()
	registerHandler(message.RequireType_RequireDisk, &diskInfo, OptionDisk)

	fileSystemInfo := b.InfoCollector.GetFileSystemInfo()
	registerHandler(message.RequireType_RequireFileSystem, &fileSystemInfo, OptionFileSystem)

	b.RegisterHandler(message.RequireType_RequirePlatform, func() (bytes []byte, option Option) {
		info := b.InfoCollector.GetEdgeNodeInfo()
		b, err := info.XXX_Marshal(nil, false)
		if err != nil {
			logger.Error(err)
		}
		return b, OptionProtoBuf | OptionPlatform
	})
}

// RegisterHandler 注册回调
func (b *HeartBeater) RegisterHandler(rt message.RequireType, handler PacketHandler) {
	b.Handlers[rt] = handler
}

// 处理连接收到的请求，直到出错或者连接中断
func (b *HeartBeater) readFromConn(c net.Conn) error {
	var err error
	r := bufio.NewReader(c)
	for {
		err = skipProtocolHead(r)
		if err != nil {
			if isConnectionBroken(err) {
				return err
			}
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "read" {
				return err
			}
			logger.Warn(err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		buf, err := readBody(r)
		if err != nil {
			if isConnectionBroken(err) {
				return err
			}
			logger.Warn(err)
			continue
		}

		var require message.RequireReport
		err = proto.Unmarshal(buf, &require)
		if err != nil {
			logger.Warn("Unmarshal Master Require Failed. ", err)
			continue
		}

		// 上报数据
		b.handleRequire(require.Requires...)
	}
}

// skipProtocolHead 不断读取数据，直到读到协议头
func skipProtocolHead(c io.Reader) error {
	b := make([]byte, 1)
	headIndex := 0
	headLen := len(protocolHead)
	for headIndex < headLen {
		n, err := c.Read(b)
		if err != nil {
			return err
		}

		if n <= 0 || b[0] != protocolHead[headIndex] {
			headIndex = 0
		} else {
			headIndex += 1
		}
	}
	return nil
}

func readBody(r *bufio.Reader) ([]byte, error) {
	lengthBuf := make([]byte, 2)
	_, err := r.Read(lengthBuf)
	if err != nil {
		return nil, err
	}

	size := int(binary.LittleEndian.Uint16(lengthBuf))
	opt, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if (Option(opt) & OptionProtoBuf) == 0 {
		// TODO: 服务端序列化模式不是 ProtoBuf
		logger.Warn("Server mode is not Protobuf")
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(r, buf)
	return buf, err
}

func (b *HeartBeater) handleRequire(requires ...message.RequireType) {
	for _, req := range requires {
		handler := b.Handlers[req]
		if handler != nil {
			p, opt := handler()
			b.sendPacket(p, opt)
		}
	}
}
