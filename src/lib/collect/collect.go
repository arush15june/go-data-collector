// Collect data from system.
//
// Exports
// -------
// GetSystemInfo map[string]string: Get a map of system information,
// 								 pack to JSON, other formats.

package collect

import (
	"encoding/json"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// Container for collected system info.

type SystemInfoStat struct {
	TotalMemory     uint64
	AvailableMemory uint64
	UsedMemory      uint64
	TotalDisk       uint64
	FreeDisk        uint64
	UsedDisk        uint64
	DiskPath        string
	Hostname        string
	OS              string
	Timestamp       uint64
}

func (sysinfo SystemInfoStat) String() string {
	val, err := json.Marshal(sysinfo)
	if err != nil {
		os.Exit(1)
	}
	return string(val)
}

const (
	// DefaultWindowsPath is the default path for disk usage stats on windows
	DefaultWindowsPath = "C:\\"
	// DefaultLinuxPath is the default path for disk usage stats on linux.
	DefaultLinuxPath = "/"
)

// getMemory retuns memory usage of the system.
func getMemory() *mem.VirtualMemoryStat {
	vmem, _ := mem.VirtualMemory()
	return vmem
}

// getDiskUsage returns disk usage.
func getDiskUsage(path string) *disk.UsageStat {
	diskusage, _ := disk.Usage(path)
	return diskusage
}

// getTemperature return system temperature sensor info.
func getTemperature() []host.TemperatureStat {
	temp, _ := host.SensorsTemperatures()
	return temp
}

// getHostInfo returns host specfici information.
func getHostInfo() *host.InfoStat {
	info, _ := host.Info()
	return info
}

func getDiskPath() string {
	if runtime.GOOS == "windows" {
		return DefaultWindowsPath
	}
	return DefaultLinuxPath
}

func GetSystemInfo() SystemInfoStat {
	diskUsagePath := getDiskPath()

	var (
		diskUsage = getDiskUsage(diskUsagePath)
		memUsage  = getMemory()
		hostInfo  = getHostInfo()
		timestamp = uint64(time.Now().Unix())
	)

	systemInfo := SystemInfoStat{
		memUsage.Total,
		memUsage.Available,
		memUsage.Used,
		diskUsage.Total,
		diskUsage.Free,
		diskUsage.Used,
		diskUsage.Path,
		hostInfo.Hostname,
		hostInfo.OS,
		timestamp,
	}

	return systemInfo
}

// GetSystemInfoEx collects system info and returns it as a map[string]string.
func GetSystemInfoEx() map[string]string {
	diskUsagePath := getDiskPath()

	var (
		diskUsage = getDiskUsage(diskUsagePath)
		memUsage  = getMemory()
		hostInfo  = getHostInfo()
		timestamp = uint64(time.Now().Unix())
	)

	systemInfo := map[string]string{
		"totalMem":  strconv.FormatUint(memUsage.Total, 10),
		"availMem":  strconv.FormatUint(memUsage.Available, 10),
		"usedMem":   strconv.FormatUint(memUsage.Used, 10),
		"totalDisk": strconv.FormatUint(diskUsage.Total, 10),
		"freeDisk":  strconv.FormatUint(diskUsage.Free, 10),
		"usedDisk":  strconv.FormatUint(diskUsage.Used, 10),
		"diskPath":  diskUsage.Path,
		"hostname":  hostInfo.Hostname,
		"os":        hostInfo.OS,
		"timestamp": strconv.FormatUint(timestamp, 10),
	}

	return systemInfo
}
