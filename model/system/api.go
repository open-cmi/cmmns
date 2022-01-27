package system

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/model"
	sysmod "github.com/open-cmi/cmmns/module/system"
	"github.com/open-cmi/cmmns/storage/db"
	"github.com/open-cmi/goutils/devutil"

	_ "github.com/lib/pq" //
)

type Option struct {
	model.Option
}

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

	sqldb := db.DB
	_, err := sqldb.Exec(sqlquery)
	if err != nil {
	}
	return
}

func Get(mo *Option, field string, value string) *Model {
	columns := model.GetColumn(Model{}, []string{})

	queryClause := fmt.Sprintf(`select %s from system_status where %s=$1`, strings.Join(columns, ","), field)
	dbsql := db.GetDB()
	row := dbsql.QueryRowx(queryClause, value)

	var mdl Model
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil
	}

	return &mdl
}

// List list
func List(option *Option) (int, []Model, error) {
	dbsql := db.GetDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from system_status")
	whereClause, args := model.BuildWhereClause(&option.Option)
	countClause += whereClause
	row := dbsql.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	columns := model.GetColumn(Model{}, []string{})
	queryClause := fmt.Sprintf(`select %s from system_status`, strings.Join(columns, ","))
	finalClause := model.BuildFinalClause(&option.Option)
	queryClause += (whereClause + finalClause)
	rows, err := dbsql.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		logger.Logger.Error(err.Error())
		return count, results, nil
	}

	for rows.Next() {
		var item Model
		err := rows.StructScan(&item)
		if err != nil {
			logger.Logger.Error(err.Error())
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

var LocalModel *Model

// StartMonitor start to Update device resource
func StartMonitor() {
	devID := devutil.GetDeviceID()
	if devID == "" {
		return
	}

	if LocalModel == nil {
		LocalModel = Get(&Option{}, "id", devID)
		if LocalModel == nil {
			LocalModel = New()
			LocalModel.ID = devID
			LocalModel.IsMaster = true
		}
	}

	LocalModel.UpdatedTime = time.Now().Unix()
	LocalModel.CPUUsage = sysmod.CPUSummary()

	diskUsed, diskTotal, diskUsedPercent := sysmod.DiskSummary()
	LocalModel.DiskTotal = diskTotal
	LocalModel.DiskUsed = diskUsed
	LocalModel.DiskUsedPercent = diskUsedPercent

	memUsed, memTotal, memUsedPercent := sysmod.MemSummary()
	LocalModel.MemTotal = memTotal
	LocalModel.MemUsed = memUsed
	LocalModel.MemUsedPercent = memUsedPercent

	netSent, netRecv := sysmod.NetRateSummary()
	LocalModel.NetRecv = netRecv
	LocalModel.NetSent = netSent

	load1, load5, load15 := sysmod.LoadSummary()
	LocalModel.LoadAvg1 = load1
	LocalModel.LoadAvg5 = load5
	LocalModel.LoadAvg15 = load15

	LocalModel.Save()

	return
}
