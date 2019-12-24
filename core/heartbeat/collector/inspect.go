package collector

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"jxcore/core/heartbeat/message"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gitlab.jiangxingai.com/applications/base-modules/internal-sdk/go-utils/logger"

	"github.com/shirou/gopsutil/disk"
)

// PipPackage pip 包
type PipPackage struct {
	Name    string `json:"name"`    // 包名
	Version string `json:"version"` // 版本号
	// {"name": "apt-clone", "version": "0.2.1"}
}

// ListPipPackages List installed packages, including editables.
func ListPipPackages(python string) map[string]string {
	if python == "" {
		python = "python3"
	}
	cmd := exec.Command(python, "-m", "pip", "list", "--format", "json")
	var packages []PipPackage
	r := make(map[string]string)

	buf, err := cmd.Output()
	if err == nil {
		if json.Unmarshal(buf, &packages) == nil {
			for _, p := range packages {
				r[p.Name] = p.Version
			}
		}
	} else {
		// pip 版本较低， Fallback 方案
		cmd = exec.Command(python, "-m", "pip", "freeze")

		buf, err = cmd.Output()
		if err != nil {
			logger.Error(err)
		} else {
			scanner := bufio.NewScanner(bytes.NewBuffer(buf))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "==") {
					slice := strings.SplitN(line, "==", 2)
					r[slice[0]] = slice[1]
				}
			}
		}
	}
	return r
}

// IsRknnInstalled 判断是否安装rknn库
func IsRknnInstalled(python string) bool {
	packages := ListPipPackages(python)

	_, ok := packages["rknn"]
	return ok
}

// ProcessorInfo 逻辑处理器信息
type ProcessorInfo struct {
	Index    int
	BogoMIPS float64
	Features []string

	CPUImplementer  int
	CPUArchitecture int
	CPUVariant      int
	CPUPart         int
	CPURevision     int
}

// RawCPUInfo CPU 信息
type RawCPUInfo struct {
	Processors []*ProcessorInfo
	Serial     string
}

// GetRawCPUInfo 获取系统 CPU 信息
func GetRawCPUInfo() RawCPUInfo {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		logger.Error(err)
		return RawCPUInfo{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	return cpuInfoFromScanner(scanner)
}

func cpuInfoFromScanner(scanner *bufio.Scanner) RawCPUInfo {

	var processors []*ProcessorInfo
	var p *ProcessorInfo
	var info RawCPUInfo

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.ContainsRune(line, ':') {
			continue
		}

		if strings.HasPrefix(line, "processor") {
			p = &ProcessorInfo{}
			processors = append(processors, p)
		}
		if p == nil {
			continue
		}

		segments := strings.Split(line, ":")
		tag := strings.TrimSpace(segments[0])
		value := strings.TrimSpace(segments[1])
		switch tag {
		case "processor":
			p.Index, _ = strconv.Atoi(value)
		case "BogoMIPS":
			p.BogoMIPS, _ = strconv.ParseFloat(value, 64)
		case "Features":
			p.Features = strings.Split(value, " ")
		case "CPU implementer":
			p.CPUImplementer, _ = hexToInt(value)
		case "CPU architecture":
			p.CPUArchitecture, _ = strconv.Atoi(value)
		case "CPU variant":
			p.CPUVariant, _ = hexToInt(value)
		case "CPU part":
			p.CPUPart, _ = hexToInt(value)
		case "CPU revision":
			p.CPURevision, _ = strconv.Atoi(value)
		case "Serial":
			info.Serial = value
		}
	}
	info.Processors = processors

	return info
}

// RawNetworkInterface 网卡信息
type RawNetworkInterface struct {
	net.Interface

	IPv4 net.IP
	IPv6 net.IP
}

// GetEdgeNodeNetworkInfo 获取网卡信息
func GetEdgeNodeNetworkInfo(skipVirtual bool) []RawNetworkInterface {
	var networks []RawNetworkInterface
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Error(err)
		return networks
	}
	for _, iface := range ifaces {
		if skipVirtual {
			// 跳过虚拟网卡
			name := iface.Name

			if isVirtualNetwork(name) || len(iface.HardwareAddr) == 0 {
				continue
			}
		}

		addrs, err := iface.Addrs()
		if err != nil {
			logger.Error(err)
			continue
		}
		network := RawNetworkInterface{
			Interface: iface,
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				ip := ipNet.IP
				mask := ip.DefaultMask()
				switch len(mask) {
				case net.IPv4len:
					ipv4 := ip.To4()
					if ipv4 != nil {
						network.IPv4 = ipv4
					}
				case 0:
					fallthrough
				case net.IPv6len:
					ipv6 := ip.To16()
					if ipv6 != nil {
						network.IPv6 = ipv6
					}
				default:
					logger.Debug("Unknown Mask: ", mask.String())
				}
			}
		}

		networks = append(networks, network)
	}
	return networks
}

// RawMemoryInfo 内存信息
type RawMemoryInfo struct {
	Total     int64
	Free      int64
	Available int64
	Extras    map[string]int64
}

// GetRawMemoryInfo 获取内存信息
func GetRawMemoryInfo() RawMemoryInfo {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		logger.Error(err)
		return RawMemoryInfo{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	return memoryInfoFromScanner(scanner)
}

func memoryInfoFromScanner(scanner *bufio.Scanner) RawMemoryInfo {
	var info RawMemoryInfo
	info.Extras = make(map[string]int64)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.ContainsRune(line, ':') {
			continue
		}

		segments := strings.Split(line, ":")
		tag := strings.TrimSpace(segments[0])
		value := strings.TrimSpace(segments[1])
		switch tag {
		case "MemTotal":
			info.Total = translateFromHuman(value)
		case "MemFree":
			info.Free = translateFromHuman(value)
		case "MemAvailable":
			info.Available = translateFromHuman(value)
		default:
			info.Extras[tag] = translateFromHuman(value)
		}
	}

	return info
}

// RawDiskInfo 磁盘信息
type RawDiskInfo struct {
	Name       string
	Type       int
	Size       int // 扇区数
	SectorSize int // 扇区大小（Bytes）
}

// GetDiskInfos 获取磁盘信息
func GetDiskInfos() []RawDiskInfo {
	fInfos, err := ioutil.ReadDir("/sys/block")
	if err != nil {
		logger.Error(err)
		return nil
	}

	var disks []RawDiskInfo
	for _, it := range fInfos {
		name := it.Name()
		if strings.HasPrefix(name, "loop") || strings.Contains(name, "ram") {
			continue
		}

		rotational := filepath.Join("/sys/block", name, "queue/rotational")
		size := filepath.Join("/sys/block", name, "size")
		sectorSize := filepath.Join("/sys/block", name, "queue/hw_sector_size")

		buf, err := ioutil.ReadFile(rotational)
		if err != nil || len(buf) == 0 {
			continue
		}

		buf1, err := ioutil.ReadFile(size)
		if err != nil {
			continue
		}

		buf2, err := ioutil.ReadFile(sectorSize)
		if err != nil {
			continue
		}

		diskInfo := RawDiskInfo{
			Name:       name,
			Type:       mustInt(string(buf[0])),
			Size:       mustInt(string(buf1)),
			SectorSize: mustInt(string(buf2)),
		}

		if diskInfo.Type == int(message.DiskType_SSD) && strings.HasPrefix(name, "mmcblk") {
			diskInfo.Type = int(message.DiskType_EMMC)
		}

		disks = append(disks, diskInfo)
	}

	return disks
}

// RawFileSystemInfo 文件系统信息
type RawFileSystemInfo struct {
	MountSegment
	*disk.UsageStat
}

// GetFileSystemInfos 获取文件系统信息
func GetFileSystemInfos() []RawFileSystemInfo {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	return filesystemInfoFromScanner(scanner)
}

// MountSegment 系统挂载信息
type MountSegment struct {
	FileSystem string
	MountPoint string
	Type       string
	Options    []string
	Dump       int
	Pass       int
}

func parseMountsFile(scanner *bufio.Scanner) []MountSegment {
	var result []MountSegment

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}

		segments := strings.Split(line, " ")
		if len(segments) != 6 {
			logger.Debug("Unknown Line: ", line)
			continue
		}
		fs := MountSegment{
			FileSystem: segments[0],
			MountPoint: segments[1],
			Type:       segments[2],
			Options:    strings.Split(segments[3], ","),
			Dump:       mustInt(segments[4]),
			Pass:       mustInt(segments[5]),
		}

		result = append(result, fs)
	}

	return result
}

func filesystemInfoFromScanner(scanner *bufio.Scanner) []RawFileSystemInfo {
	mounts := parseMountsFile(scanner)
	var result []RawFileSystemInfo

	for index := 0; index < len(mounts); index++ {
		it := mounts[index]
		if !strings.HasPrefix(it.FileSystem, "/") {
			continue
		}
		if strings.HasPrefix(it.MountPoint, "/var/lib/docker/volumes/") {
			continue
		}

		var fsInfo RawFileSystemInfo
		fsInfo.MountSegment = it

		du, err := disk.Usage(it.MountPoint)
		if err != nil {
			logger.Error(err)
		} else {
			fsInfo.UsageStat = du
		}

		result = append(result, fsInfo)
	}

	return result
}
