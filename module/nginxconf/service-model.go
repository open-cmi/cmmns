package nginxconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
)

const ServiceModelKey = "nginx-service-port"

type ServiceModel struct {
	HttpPort  int  `json:"http_port"`
	HttpsPort int  `json:"https_port"`
	isNew     bool `json:"-"`
}

func (m *ServiceModel) Key() string {
	return ServiceModelKey
}

func (m *ServiceModel) Value() string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (m *ServiceModel) Save() error {
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

func GetServiceModel() *ServiceModel {
	queryClause := "select value from k_v_table where key=$1"
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, ServiceModelKey)
	if row == nil {
		return nil
	}
	var v string
	err := row.Scan(&v)
	if err != nil {
		logger.Errorf("nginx service model row scan failed: %s\n", err.Error())
		return nil
	}
	var m ServiceModel
	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		logger.Errorf("nginx service model model unmarshal failed: %s\n", err.Error())
		return nil
	}
	return &m
}

func NewServiceModel() *ServiceModel {
	return &ServiceModel{
		isNew: true,
	}
}
