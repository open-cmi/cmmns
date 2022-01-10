package system

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sysmod "github.com/open-cmi/cmmns/modules/system"
	"github.com/open-cmi/cmmns/storage/db"
	"github.com/open-cmi/goutils/device"

	_ "github.com/lib/pq" //
)

// SystemInfo device info
type SystemInfo struct {
	DeviceID        string  `json:"deviceid"`
	CPUUsage        float64 `json:"cpu_usage"`
	DiskUsed        uint64  `json:"disk_used"`
	DiskTotal       uint64  `json:"disk_total"`
	DiskUsedPercent float64 `json:"disk_used_percent"`
	MemUsed         uint64  `json:"mem_used"`
	MemTotal        uint64  `json:"mem_total"`
	MemUsedPercent  float64 `json:"mem_used_percent"`
	NetSent         uint64  `json:"net_sent"`
	NetRecv         uint64  `json:"net_recv"`
	LoadAvg1        float64 `json:"load_avg_1"`
	LoadAvg5        float64 `json:"load_avg_5"`
	LoadAvg15       float64 `json:"load_avg_15"`
}

// GetSystemInfo get device info
func GetSystemInfo() (info SystemInfo, err error) {
	columns := []string{
		"deviceid",
		"cpu_usage",
		"disk_used",
		"disk_total",
		"disk_used_percent",
		"mem_used",
		"mem_total",
		"mem_used_percent",
		"net_sent",
		"net_recv",
		"load_avg_1",
		"load_avg_5",
		"load_avg_15",
	}
	dbquery := fmt.Sprintf("select %s from systeminfo", strings.Join(columns, ","))
	sqldb := db.DB
	row := sqldb.QueryRow(dbquery)
	if row == nil {
		return info, errors.New("get db system info failed")
	}

	err = row.Scan(&info.DeviceID, &info.CPUUsage, &info.DiskUsed, &info.DiskTotal, &info.DiskUsedPercent,
		&info.MemUsed, &info.MemTotal, &info.MemUsedPercent, &info.NetSent, &info.NetRecv, &info.LoadAvg1,
		&info.LoadAvg5, &info.LoadAvg15)
	if err != nil {
		return info, err
	}
	return info, nil
}

// InitDeviceID init device id
func InitDeviceID() {
	utime := time.Now().UTC().Format(time.RFC3339)
	deviceid := device.GetDeviceID()
	if deviceid == "" {
		return
	}

	sqlquery := fmt.Sprintf(`insert into systeminfo 
		(utime, deviceid, is_master) 
		VALUES ('%s', '%s', '%t') 
		ON CONFLICT (deviceid) DO
		NOTHING`,
		utime, deviceid, true)

	sqldb := db.DB
	_, err := sqldb.Exec(sqlquery)
	if err != nil {
	}
	return
}

// UpdateSystemInfo update device info
func UpdateSystemInfo() {
	deviceid := device.GetDeviceID()
	if deviceid == "" {
		return
	}

	utime := time.Now().UTC().Format(time.RFC3339)

	cpuUsage := sysmod.CPUSummary()

	diskUsed, diskTotal, diskUsedPercent := sysmod.DiskSummary()

	memUsed, memTotal, memUsedPercent := sysmod.MemSummary()

	netSent, netRecv := sysmod.NetRateSummary()

	load1, load5, load15 := sysmod.LoadSummary()

	sqlquery := fmt.Sprintf(`update systeminfo set utime='%s',
		cpu_usage=%f, disk_total=%d, disk_used=%d, disk_used_percent=%f,
		mem_total=%d, mem_used=%d, mem_used_percent=%f, 
		load_avg_1=%f, load_avg_5=%f, load_avg_15=%f,
		net_sent=%d, net_recv=%d where deviceid='%s'
		`, utime, cpuUsage, diskTotal, diskUsed, diskUsedPercent,
		memTotal, memUsed, memUsedPercent, load1, load5, load15,
		netSent, netRecv, deviceid)

	sqldb := db.DB
	_, err := sqldb.Exec(sqlquery)
	if err != nil {
	}
	return
}

// StartMonitor start to UpdateSystemInfo device resource
func StartMonitor() {
	UpdateSystemInfo()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			UpdateSystemInfo()
		}
	}
}
