package time

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

type Setting struct {
	TimeZone   string `json:"timezone"`
	NtpServer  string `json:"ntp_server"`
	AutoAdjust bool   `json:"auto_adjust"`
	isNew      bool
}

func (s *Setting) Key() string {
	return "time-setting"
}

func (s *Setting) Value() string {
	v, _ := json.Marshal(s)
	return string(v)
}

func (m *Setting) Save() error {
	db := sqldb.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := []string{"key", "value"}
		values := []string{"$1", "$2"}

		insertClause := fmt.Sprintf("insert into k_v_table(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)
		_, err := db.Exec(insertClause, m.Key(), m.Value())
		if err != nil {
			logger.Errorf("create setting failed: %s", err.Error())
			return errors.New("create setting failed")
		}
		m.isNew = false
	} else {
		updateClause := "update k_v_table set value=$1 where key=$2"
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.Exec(updateClause, m.Value(), m.Key())
		if err != nil {
			logger.Errorf("update setting failed: %s", err.Error())
			return errors.New("update setting failed")
		}
	}

	return nil
}

func New() *Setting {
	return &Setting{
		isNew: true,
	}
}

func Get() *Setting {
	db := sqldb.GetDB()

	var m Setting
	queryClause := `select value from k_v_table where key=$1`
	row := db.QueryRowx(queryClause, m.Key())
	if row == nil {
		return nil
	}
	var value string
	err := row.Scan(&value)
	if err != nil {
		logger.Infof("ntp setting scan setting failed: %s\n", err.Error())
		return nil
	}

	err = json.Unmarshal([]byte(value), &m)
	if err != nil {
		logger.Errorf("ntp setting json unmarshal failed: %s\n", err.Error())
		return nil
	}

	return &m
}
