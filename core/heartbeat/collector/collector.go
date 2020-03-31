package collector

import (
	"github.com/JK-97/edge-guard/core/heartbeat/message"
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
