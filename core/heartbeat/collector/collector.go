package collector

import (
	"jxcore/core/heartbeat/message"
)

type InfoCollector interface {
	GetWorkerID() message.WorkerID
	GetCPUInfo() message.CPUInfo
	GetNetworkInfo() message.NetworkInfo
	GetDiskStat() message.DiskInfo
	GetFileSystemInfo() message.FileSystemInfo
	GetMemoryInfo() message.MemoryInfo
	GetEdgeNodeInfo() message.PlatformInfo
}
