package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
)

type DevModel struct {
	Dev          string `json:"dev"`
	DHCP         bool   `json:"dhcp"`
	Address      string `json:"address,omitempty"`
	Netmask      string `json:"netmask,omitempty"`
	Gateway      string `json:"gateway,omitempty"`
	PreferredDNS string `json:"preferred_dns,omitempty"`
	AlternateDNS string `json:"alternate_dns,omitempty"`
	isNew        bool   `json:"-"`
}

func (m *DevModel) Key() string {
	return fmt.Sprintf("net-config-%s", m.Dev)
}

func (m *DevModel) Value() string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (m *DevModel) Save() error {
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

func (m *DevModel) Remove() error {
	deleteClause := "delete from k_v_table where key=$1"
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, m.Key())
	if err != nil {
		return errors.New("del net config failed")
	}
	return err
}

func New() *DevModel {
	return &DevModel{
		isNew: true,
	}
}

func GetNetConfig(dev string) *DevModel {
	queryClause := "select value from k_v_table where key=$1"
	db := sqldb.GetDB()
	key := fmt.Sprintf("net-config-%s", dev)
	row := db.QueryRowx(queryClause, key)
	if row == nil {
		return nil
	}
	var v string
	err := row.Scan(&v)
	if err != nil {
		logger.Errorf("row scan failed: %s\n", err.Error())
		return nil
	}
	var m DevModel
	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		logger.Errorf("DevModel unmarshal failed: %s\n", err.Error())
		return nil
	}
	return &m
}
