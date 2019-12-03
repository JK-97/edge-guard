package hearbeat

import (
	"context"
	"fmt"
	"jxcore/core/device"
	"jxcore/internal/network"
	"net"
	"time"

	"gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"
)

var (
	conn         net.Conn
	connMasterip string
)

// AliveReport 上报心跳，直到连续出错次数超过 allowContinuousFailed
func AliveReport(ctx context.Context, masterip string, allowContinuousFailed int) error {
	msg, err := getHeartbeatMsg()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(heartBeatInterval)
	defer ticker.Stop()

	continuousFailed := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if tryToSend(masterip, msg) != nil {
				continuousFailed++
				logger.Infof("Heartbeat to %s failed %d times.", masterip, continuousFailed)
				if continuousFailed >= allowContinuousFailed {
					return nil
				}
			} else {
				continuousFailed = 0
			}
		}
	}
}

func getHeartbeatMsg() (string, error) {
	currentdevice, err := device.GetDevice()
	if err != nil {
		return "", err
	}
	vpn, err := currentdevice.Vpn.Interface()
	if err != nil {
		return "", err
	}
	vpninterface, err := network.GetMyIP(vpn)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(vpninterface + ":" + currentdevice.WorkerID), nil
}

// 尝试发送 msg 到 masterip，如果masterip不变，会复用上次的连接
func tryToSend(masterip, msg string) (err error) {
	if conn != nil && masterip != connMasterip {
		conn.Close()
		conn = nil
	}

	if conn == nil {
		conn, err = net.DialTimeout("tcp", masterip+":30431", dialTimeout)
		if err != nil {
			return err
		}
		connMasterip = masterip
	}

	if _, err = conn.Write([]byte(msg)); err != nil {
		conn.Close()
		conn = nil
		return err
	}
	return nil
}
