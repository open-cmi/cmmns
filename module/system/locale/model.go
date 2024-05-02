package locale

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

const LocaleKey = "locale"

type LocaleModel struct {
	Lang  string `json:"lang" db:"lang"`
	isNew bool   `json:"-"`
}

func (em *LocaleModel) Key() string {
	return LocaleKey
}

func (em *LocaleModel) Value() string {
	v, _ := json.Marshal(em)
	return string(v)
}

func (m *LocaleModel) Save() error {
	db := sqldb.GetConfDB()

	if m.isNew {
		// 存储到数据库
		columns := []string{"key", "value"}
		values := []string{"$1", "$2"}

		insertClause := fmt.Sprintf("insert into k_v_table(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)
		_, err := db.Exec(insertClause, m.Key(), m.Value())
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
		m.isNew = false
	} else {
		updateClause := "update k_v_table set value=$1 where key=$2"
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.Exec(updateClause, m.Value(), m.Key())
		if err != nil {
			logger.Errorf("update model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}
	return nil
}

func New() *LocaleModel {
	return &LocaleModel{
		isNew: true,
	}
}

func Get() *LocaleModel {
	queryClause := "select value from k_v_table where key=$1"
	db := sqldb.GetConfDB()
	row := db.QueryRowx(queryClause, LocaleKey)
	if row == nil {
		return nil
	}
	var v string
	err := row.Scan(&v)
	if err != nil {
		logger.Errorf("locale row scan failed: %s\n", err.Error())
		return nil
	}
	var m LocaleModel
	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		logger.Errorf("local model unmarshal failed: %s\n", err.Error())
		return nil
	}
	return &m
}
