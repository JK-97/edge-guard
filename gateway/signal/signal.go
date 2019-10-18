package signal

import "os"

// StopSignals 停止信号列表
var StopSignals []os.Signal

// HangUpSignals 父进程终止信号列表
var HangUpSignals []os.Signal
