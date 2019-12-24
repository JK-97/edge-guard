package heartbeat

import (
	"context"
	"errors"
	"net"
)

// ErrMaxErrorCountExceed 出错次数达到最大值
var ErrMaxErrorCountExceed = errors.New("max error count exceed")

const (
	v2HeartbeatPort = "31431" // v2 版心跳协议，服务端监听端口
)

// AliveReport 上报心跳，直到连续出错次数超过 allowContinuousFailed
func AliveReport(ctx context.Context, masterip string, allowContinuousFailed int) error {
	beater := NewDefaultHeartBeater(net.JoinHostPort(masterip, v2HeartbeatPort), allowContinuousFailed)
	return beater.Beat(ctx)
}
