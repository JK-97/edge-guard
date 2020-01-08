// +build linux

package heartbeat

import "syscall"

// 需要重点关注的异常
var (
	AbortErrs = []error{
		syscall.Errno(32), // broken pipe
	} // 连接断开服务恢复的异常
	ConnectErrs = []error{
		syscall.Errno(111),
	} // 连接异常
)
