package heartbeat

import (
	"net"
	"time"
)

// Connection 心跳连接
type Connection struct {
	BeatPack []byte        // 心跳包
	Span     time.Duration // 心跳间隔
	Conn     net.Conn      // 连接
}

// Send 发送心跳包
func (c *Connection) Send() bool {
	_, err := c.Conn.Write(c.BeatPack)

	return err == nil
}
