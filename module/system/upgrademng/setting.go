package upgrademng

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

type UpgradeSettingModel struct {
	AutoUpgrade bool `json:"auto_upgrade"`
	isNew       bool `json:"-"`
}

func (m *UpgradeSettingModel) Key() string {
	return "system-upgrade-setting"
}

func (m *UpgradeSettingModel) Value() string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (m *UpgradeSettingModel) Save() error {
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

func NewUpgradeSettingModel() *UpgradeSettingModel {
	return &UpgradeSettingModel{
		isNew: true,
	}
}

func GetUpgradeSettingModel() *UpgradeSettingModel {
	queryClause := "select value from k_v_table where key=$1"
	db := sqldb.GetDB()

	var mdl UpgradeSettingModel
	row := db.QueryRowx(queryClause, mdl.Key())
	if row == nil {
		return nil
	}
	var v string
	err := row.Scan(&v)
	if err != nil {
		logger.Errorf("upgrade row scan failed: %s\n", err.Error())
		return nil
	}
	err = json.Unmarshal([]byte(v), &mdl)
	if err != nil {
		logger.Errorf("upgrade unmarshal failed: %s\n", err.Error())
		return nil
	}
	return &mdl
}

type SetSettingRequest struct {
	AutoUpgrade bool `json:"auto_upgrade"`
}

func SetUpgradeSetting(req *SetSettingRequest) error {
	m := GetUpgradeSettingModel()
	if m == nil {
		m = NewUpgradeSettingModel()
	}
	m.AutoUpgrade = req.AutoUpgrade
	return m.Save()
}

func GetUpgradeSetting() *UpgradeSettingModel {
	m := GetUpgradeSettingModel()
	if m == nil {
		m = NewUpgradeSettingModel()
		m.AutoUpgrade = false
	}
	return m
}
