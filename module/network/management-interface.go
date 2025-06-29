package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
)

type ManagementInterfaceModel struct {
	Interfaces []string `json:"interfaces"`
	isNew      bool     `json:"-"`
}

func (m *ManagementInterfaceModel) Key() string {
	return "network-management-interface"
}

func (m *ManagementInterfaceModel) Value() string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (m *ManagementInterfaceModel) Save() error {
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

func (m *ManagementInterfaceModel) Remove() error {
	deleteClause := "delete from k_v_table where key=$1"
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, m.Key())
	if err != nil {
		return errors.New("delete management interface failed")
	}
	return err
}

func GetManagementInterfaceModel() *ManagementInterfaceModel {
	var m ManagementInterfaceModel
	queryClause := "select value from k_v_table where key=$1"
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, m.Key())
	if row == nil {
		return nil
	}
	var v string
	err := row.Scan(&v)
	if err != nil {
		logger.Errorf("row scan failed: %s\n", err.Error())
		return nil
	}
	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		logger.Errorf("management interface unmarshal failed: %s\n", err.Error())
		return nil
	}
	return &m
}

func LoadNetworkManagementInterface() []string {
	m := GetManagementInterfaceModel()
	if m == nil {
		return []string{}
	}
	return m.Interfaces
}
