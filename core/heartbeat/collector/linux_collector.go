package collector

import (
	"github.com/JK-97/edge-guard/core/device"
	"github.com/JK-97/edge-guard/core/heartbeat/message"
)

type linuxCollector struct{}

func NewLinuxCollector() InfoCollector {
	return &linuxCollector{}
}

func (l *linuxCollector) GetWorkerID() message.WorkerID {
	currentDevice, _ := device.GetDevice()
	return message.WorkerID{WorkerId: currentDevice.WorkerID}
}

func (l *linuxCollector) GetCPUInfo() message.CPUInfo {
	var info message.CPUInfo
	raw := GetRawCPUInfo()
	info.Cores = uint32(len(raw.Processors))
	return info
}

func (l *linuxCollector) GetDiskStat() message.DiskInfo {
	var info message.DiskInfo
	for _, it := range GetDiskInfos() {
		stat := message.DiskStat{
			Name:     it.Name,
			Type:     message.DiskType(it.Type),
			Capacity: uint64(it.Size) * uint64(it.SectorSize),
		}

		info.Disks = append(info.Disks, &stat)
	}

	return info
}

func (l *linuxCollector) GetFileSystemInfo() message.FileSystemInfo {
	var info message.FileSystemInfo

	for _, it := range GetFileSystemInfos() {
		stat := message.FileSystem{
			FilePath:  it.FileSystem,
			MountPath: it.MountPoint,
			Format:    it.Type,
		}
		if it.UsageStat != nil {
			stat.Capacity = it.Total
		}
		info.FileSystems = append(info.FileSystems, &stat)
	}

	return info
}

func (l *linuxCollector) GetMemoryInfo() message.MemoryInfo {
	return message.MemoryInfo{
		Capacity: uint64(GetRawMemoryInfo().Total),
	}
}

func (l *linuxCollector) GetNetworkInfo() message.NetworkInfo {
	var info message.NetworkInfo
	var networks []*message.NetworkInterface

	ifaces := GetEdgeNodeNetworkInfo(true)

	for _, iface := range ifaces {
		network := message.NetworkInterface{
			Name: iface.Name,
			Ipv4: iface.IPv4,
			Ipv6: iface.IPv6,
			Mac:  []byte(iface.HardwareAddr),
		}

		networks = append(networks, &network)
	}
	info.Networks = networks

	return info
}

func (l *linuxCollector) GetEdgeNodeInfo() message.PlatformInfo {
	return message.PlatformInfo{
		AiPlatform: &message.AIPlatform{
			IncludeRknn: IsRknnInstalled(""),
		},
	}
}
