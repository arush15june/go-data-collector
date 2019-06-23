// Collect data from system.
//
// Exports
// -------
// GetSystemInfo map[string]string: Get a map of system information,
// 								 pack to JSON, other formats.

package collect

import (
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// DefaultWindowsPath is the default path for disk usage stats on windows
const DefaultWindowsPath = "C:\\"

// DefaultLinuxPath is the default path for disk usage stats on linux.
const DefaultLinuxPath = "/"

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

// GetSystemInfo collects system info and returns it as a map[string]string.
func GetSystemInfo() map[string]string {
	var diskUsagePath string
	if runtime.GOOS == "windows" {
		diskUsagePath = DefaultWindowsPath
	} else {
		diskUsagePath = DefaultLinuxPath
	}

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
