package system

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/pkg/dev"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

// InitDevID init device id
func InitDevID() {
	utime := time.Now().UTC().Format(time.RFC3339)
	deviceid := dev.GetDeviceID()
	if deviceid == "" {
		return
	}

	sqlquery := fmt.Sprintf(`insert into system_status 
		(utime, deviceid, is_master) 
		VALUES ('%s', '%s', '%t') 
		ON CONFLICT (deviceid) DO
		NOTHING`,
		utime, deviceid, true)

	db := sqldb.GetDB()
	_, err := db.Exec(sqlquery)
	if err != nil {
		logger.Errorf("update system si failed: %s\n", err.Error())
	}
}

type SystemHostInfo struct {
	DevID      string `json:"dev_id"`
	Hostname   string `json:"hostname"`
	CPUCores   int    `json:"cpu_cores"`
	CPUThreads int    `json:"cpu_threads"`
}

type SystemStatusInfo struct {
	CPUUsage        float64 `json:"cpu_usage"`
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
	si.Hostname, err = os.Hostname()
	if err != nil {
		return si, err
	}

	si.CPUCores, si.CPUThreads, _ = CPUSummary()
	return si, nil
}

func GetSystemStatusInfo() (si SystemStatusInfo, err error) {

	_, _, si.CPUUsage = CPUSummary()

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

func Get(mo *goparam.Param, field string, value string) *Model {
	columns := goparam.GetColumn(Model{}, []string{})

	queryClause := fmt.Sprintf(`select %s from system_status where %s=$1`, strings.Join(columns, ","), field)
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, value)

	var mdl Model
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Errorf("row scan failed: %s\n", err.Error())
		return nil
	}

	return &mdl
}

// List list
func List(option *goparam.Param) (int, []Model, error) {
	db := sqldb.GetDB()

	var results []Model = []Model{}

	countClause := "select count(*) from system_status"
	row := db.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	columns := goparam.GetColumn(Model{}, []string{})
	queryClause := fmt.Sprintf(`select %s from system_status`, strings.Join(columns, ","))
	finalClause := goparam.BuildFinalClause(option)
	queryClause += finalClause
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		logger.Error(err.Error())
		return count, results, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item Model
		err := rows.StructScan(&item)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

var LocalModel *Model
var devID string

// StartMonitor start to Update device resource
func StartMonitor() {

	if devID == "" {
		devID = dev.GetDeviceID()
		if devID == "" {
			return
		}
	}

	if LocalModel == nil {
		LocalModel = Get(nil, "id", devID)
		if LocalModel == nil {
			LocalModel = New()
			LocalModel.ID = devID
			LocalModel.IsMaster = true
		}
	}

	LocalModel.UpdatedTime = time.Now().Unix()

	_, _, LocalModel.CPUUsage = CPUSummary()

	diskUsed, diskTotal, diskUsedPercent := DiskSummary()
	LocalModel.DiskTotal = diskTotal
	LocalModel.DiskUsed = diskUsed
	LocalModel.DiskUsedPercent = diskUsedPercent

	memUsed, memTotal, memUsedPercent := MemSummary()
	LocalModel.MemTotal = memTotal
	LocalModel.MemUsed = memUsed
	LocalModel.MemUsedPercent = memUsedPercent

	netSent, netRecv := NetRateSummary()
	LocalModel.NetRecv = netRecv
	LocalModel.NetSent = netSent

	load1, load5, load15 := LoadSummary()
	LocalModel.LoadAvg1 = load1
	LocalModel.LoadAvg5 = load5
	LocalModel.LoadAvg15 = load15

	LocalModel.Save()
	LocalModel.IsNew = false
}
