package system

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/model"
	"github.com/open-cmi/cmmns/storage/db"
)

// Model system model
type Model struct {
	ID              string  `json:"id" db:"id"`
	UpdatedTime     int64   `json:"updated_time" db:"updated_time"`
	IsMaster        bool    `json:"is_master" db:"is_master"`
	CPUUsage        float64 `json:"cpu_usage" db:"cpu_usage"`
	DiskUsed        uint64  `json:"disk_used" db:"disk_used"`
	DiskTotal       uint64  `json:"disk_total" db:"disk_total"`
	DiskUsedPercent float64 `json:"disk_used_percent" db:"disk_used_percent"`
	MemUsed         uint64  `json:"mem_used" db:"mem_used"`
	MemTotal        uint64  `json:"mem_total" db:"mem_total"`
	MemUsedPercent  float64 `json:"mem_used_percent" db:"mem_used_percent"`
	NetSent         uint64  `json:"net_sent" db:"net_sent"`
	NetRecv         uint64  `json:"net_recv" db:"net_recv"`
	LoadAvg1        float64 `json:"load_avg_1" db:"load_avg_1"`
	LoadAvg5        float64 `json:"load_avg_5" db:"load_avg_5"`
	LoadAvg15       float64 `json:"load_avg_15" db:"load_avg_15"`
	IsNew           bool
}

func (m *Model) Save() error {
	sqldb := db.GetDB()

	if m.IsNew {
		// 存储到数据库
		columns := model.GetColumn(*m, []string{})
		values := model.GetColumnNamed(columns)

		insertClause := fmt.Sprintf("insert into system_status(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Logger.Debug("start to exec sql clause: %s", insertClause)

		_, err := sqldb.NamedExec(insertClause, m)
		if err != nil {
			logger.Logger.Error("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
	} else {
		columns := model.GetColumn(*m, []string{"id", "created_time"})

		m.UpdatedTime = time.Now().Unix()
		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update system_status set %s where id=:id", strings.Join(updates, ","))
		logger.Logger.Debug("start to exec sql clause: %s", updateClause)
		_, err := sqldb.NamedExec(updateClause, m)
		if err != nil {
			logger.Logger.Error("update system_status model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}
	return nil
}

func New() *Model {
	return &Model{
		IsNew: true,
	}
}
