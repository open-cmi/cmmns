package system

import (
	"os"
	"runtime"

	"github.com/open-cmi/gobase/pkg/dev"
	"github.com/shirou/gopsutil/v4/host"
)

type SystemHostInfo struct {
	DevID         string  `json:"dev_id"`
	Hostname      string  `json:"hostname"`
	OS            string  `json:"os"`
	Arch          string  `json:"arch"`
	UptimeDays    uint64  `json:"uptime_days"`
	UptimeHours   uint64  `json:"uptime_hours"`
	UptimeMinutes uint64  `json:"uptime_minutes"`
	LoadAvg1      float64 `json:"load_avg_1"`
	LoadAvg5      float64 `json:"load_avg_5"`
	LoadAvg15     float64 `json:"load_avg_15"`
}

type SystemStatusInfo struct {
	CPUUsage        float64 `json:"cpu_usage"`
	CPUCores        int     `json:"cpu_cores"`
	CPUThreads      int     `json:"cpu_threads"`
	DiskUsed        uint64  `json:"disk_used"`
	DiskTotal       uint64  `json:"disk_total"`
	DiskUsedPercent float64 `json:"disk_used_percent"`
	MemUsed         uint64  `json:"mem_used"`
	MemTotal        uint64  `json:"mem_total"`
	MemUsedPercent  float64 `json:"mem_used_percent"`
	LoadAvg1        float64 `json:"load_avg_1"`
	LoadAvg5        float64 `json:"load_avg_5"`
	LoadAvg15       float64 `json:"load_avg_15"`
}

type NetLoadInfo struct {
	NetSent uint64 `json:"net_sent"`
	NetRecv uint64 `json:"net_recv"`
}

func GetNetLoadInfo() NetLoadInfo {
	netSent, netRecv := NetRateSummary()
	var ns NetLoadInfo
	ns.NetRecv = netRecv
	ns.NetSent = netSent
	return ns
}

func GetBasicHostInfo() (si SystemHostInfo, err error) {
	si.DevID = dev.GetDeviceID()
	si.OS = runtime.GOOS
	si.Arch = runtime.GOARCH
	totalSeconds, _ := host.Uptime()
	si.UptimeDays = totalSeconds / 86400
	si.UptimeHours = (totalSeconds % 86400) / 3600
	si.UptimeMinutes = (totalSeconds % 3600) / 60

	load1, load5, load15 := LoadSummary()
	si.LoadAvg1 = load1
	si.LoadAvg5 = load5
	si.LoadAvg15 = load15

	si.Hostname, err = os.Hostname()
	return si, err
}

func GetSystemStatusInfo() (si SystemStatusInfo, err error) {

	si.CPUCores, si.CPUThreads, si.CPUUsage = CPUSummary()

	diskUsed, diskTotal, diskUsedPercent := DiskSummary()
	si.DiskTotal = diskTotal
	si.DiskUsed = diskUsed
	si.DiskUsedPercent = diskUsedPercent

	memUsed, memTotal, memUsedPercent := MemSummary()
	si.MemTotal = memTotal
	si.MemUsed = memUsed
	si.MemUsedPercent = memUsedPercent

	load1, load5, load15 := LoadSummary()
	si.LoadAvg1 = load1
	si.LoadAvg5 = load5
	si.LoadAvg15 = load15
	return si, nil
}
