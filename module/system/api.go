package system

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/goutils/devutil"
)

// InitDevID init device id
func InitDevID() {
	utime := time.Now().UTC().Format(time.RFC3339)
	deviceid := devutil.GetDeviceID()
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
	}
	return
}

func Get(mo *api.Option, field string, value string) *Model {
	columns := api.GetColumn(Model{}, []string{})

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
func List(option *api.Option) (int, []Model, error) {
	db := sqldb.GetDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from system_status")
	whereClause, args := api.BuildWhereClause(option)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	columns := api.GetColumn(Model{}, []string{})
	queryClause := fmt.Sprintf(`select %s from system_status`, strings.Join(columns, ","))
	finalClause := api.BuildFinalClause(option)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		logger.Error(err.Error())
		return count, results, nil
	}

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
		devID = devutil.GetDeviceID()
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
	LocalModel.CPUUsage = CPUSummary()

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
