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

type CPUUsageModel struct {
	DevID       string  `json:"dev_id" db:"dev_id"`
	Step        int     `json:"step" db:"step"`
	UpdatedTime int64   `json:"updated_time" db:"updated_time"`
	CPUUsage    float64 `json:"cpu_usage" db:"cpu_usage"`
}

func (m *CPUUsageModel) Save() error {
	db := sqldb.GetDB()

	// 存储到数据库
	columns := goparam.GetColumn(*m, []string{})
	values := goparam.GetColumnInsertNamed(columns)

	updateColumns := goparam.GetColumn(*m, []string{})
	updateValues := goparam.GetColumnUpdateNamed(updateColumns)

	insertClause := fmt.Sprintf("insert into system_cpu_usage(%s) values(%s) on conflict(dev_id,step) do update set %s",
		strings.Join(columns, ","), strings.Join(values, ","), strings.Join(updateValues, ","))

	logger.Debugf("start to exec sql clause: %s", insertClause)

	_, err := db.NamedExec(insertClause, m)
	if err != nil {
		logger.Errorf("create or update model failed: %s", err.Error())
		return errors.New("create or update model failed")
	}

	return nil
}

var step int = -1
var round int = 60 / 15 * 60

func GetMaxStepCpuUsageModel() *CPUUsageModel {
	queryClause := "select * from system_cpu_usage order by step desc"
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause)
	if row == nil {
		return nil
	}
	var m CPUUsageModel
	err := row.StructScan(&m)
	if err != nil {
		logger.Errorf("struct scan failed: %s\n", err.Error())
		return nil
	}

	return &m
}

func MonitorCpuUsage() {
	if devID == "" {
		devID = dev.GetDeviceID()
		if devID == "" {
			return
		}
	}

	if step == -1 {
		m := GetMaxStepCpuUsageModel()
		if m != nil {
			step = m.Step
		}
	}
	step = (step + 1) % round
	n := &CPUUsageModel{}
	n.DevID = devID
	n.UpdatedTime = time.Now().Unix()
	n.Step = step
	_, _, n.CPUUsage = CPUSummary()

	err := n.Save()
	if err != nil {
		logger.Errorf("cpu usage monitor save failed: %s\n", err.Error())
	}
}
