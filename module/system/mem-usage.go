package system

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/dev"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

type MemUsageModel struct {
	DevID       string  `json:"dev_id" db:"dev_id"`
	Step        int     `json:"step" db:"step"`
	UpdatedTime int64   `json:"updated_time" db:"updated_time"`
	MemUsage    float64 `json:"mem_usage" db:"mem_usage"`
}

func (m *MemUsageModel) Save() error {
	db := sqldb.GetDB()

	// 存储到数据库
	columns := goparam.GetColumn(*m, []string{})
	values := goparam.GetColumnInsertNamed(columns)

	updateColumns := goparam.GetColumn(*m, []string{"dev_id", "step"})
	updateNames := goparam.GetColumnUpdateNamed(updateColumns)

	insertClause := fmt.Sprintf("insert into system_mem_usage(%s) values(%s) on conflict(dev_id,step) do update set %s",
		strings.Join(columns, ","), strings.Join(values, ","), strings.Join(updateNames, ","))

	logger.Debugf("start to exec sql clause: %s", insertClause)
	_, err := db.NamedExec(insertClause, m)
	if err != nil {
		logger.Errorf("create or update model failed: %s", err.Error())
		return errors.New("create or update model failed")
	}

	return nil
}

var memStep int = -1
var memRound int = 60 / 15 * 60

func GetMaxStepMemUsageModel() *MemUsageModel {
	queryClause := "select * from system_mem_usage order by step desc"
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause)
	if row == nil {
		return nil
	}
	var m MemUsageModel
	err := row.StructScan(&m)
	if err != nil {
		logger.Errorf("struct scan failed: %s\n", err.Error())
		return nil
	}

	return &m
}

func MonitorMemUsage() {
	if devID == "" {
		devID = dev.GetDeviceID()
		if devID == "" {
			return
		}
	}

	if memStep == -1 {
		m := GetMaxStepMemUsageModel()
		if m != nil {
			memStep = m.Step
		}
	}
	memStep = (memStep + 1) % memRound
	n := &MemUsageModel{}
	n.DevID = devID
	n.UpdatedTime = time.Now().Unix()
	n.Step = memStep
	_, _, n.MemUsage = MemSummary()

	err := n.Save()
	if err != nil {
		logger.Errorf("memory usage monitor save failed: %s\n", err.Error())
	}
}
