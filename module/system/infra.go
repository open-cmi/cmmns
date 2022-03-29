package system

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// NetRateSummary 网络收发
func NetRateSummary() (uint64, uint64) {
	firstNetStat, err := net.IOCounters(false)
	if err != nil || len(firstNetStat) == 0 {
		return 0, 0
	}

	time.Sleep(1 * time.Second)

	secondNetStat, err := net.IOCounters(false)
	if err != nil || len(secondNetStat) == 0 {
		return 0, 0
	}

	sendPercent := secondNetStat[0].BytesSent - firstNetStat[0].BytesSent
	receivePercent := secondNetStat[0].BytesRecv - firstNetStat[0].BytesRecv

	return sendPercent, receivePercent
}

// CPUSummary get cpu usgae
func CPUSummary() (usage float64) {
	arr, err := cpu.Percent(time.Second, false)
	if err != nil || len(arr) == 0 {
		return
	}
	usage = arr[0]
	return
}

// MemSummary memory usage
func MemSummary() (used uint64, total uint64, usedPercent float64) {
	memstat, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	total = memstat.Total
	used = memstat.Used
	usedPercent = memstat.UsedPercent
	return
}

// DiskSummary  disk summary , only stat root dir/
func DiskSummary() (used uint64, total uint64, usedPercent float64) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return
	}

	for _, part := range parts {
		if part.Mountpoint == "/" {
			diskStat, _ := disk.Usage(part.Mountpoint)
			used = diskStat.Used
			total = diskStat.Total
			usedPercent = diskStat.UsedPercent
			break
		}
	}
	return
}

// LoadSummary load summary
func LoadSummary() (load1 float64, load5 float64, load15 float64) {
	stat, err := load.Avg()
	if err != nil {
		return
	}
	return stat.Load1, stat.Load5, stat.Load15
}
