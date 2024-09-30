package prod

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/service/initial"
)

type ProdModel struct {
	Name   string `json:"name"`
	Footer string `json:"footer"`
	isNew  bool
}

func (m *ProdModel) Key() string {
	return "prod-model-info"
}

func (m *ProdModel) Value() string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (m *ProdModel) Save() error {
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

func GetProdModel() *ProdModel {
	var m ProdModel

	queryClause := "select value from k_v_table where key=$1"
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, m.Key())
	if row == nil {
		return nil
	}

	var v string
	err := row.Scan(&v)
	if err != nil {
		logger.Errorf("Prod model row scan failed: %s\n", err.Error())
		return nil
	}

	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		logger.Errorf("Prod model unmarshal failed: %s\n", err.Error())
		return nil
	}
	return &m
}
func NewProdModel() *ProdModel {
	return &ProdModel{
		isNew: true,
	}
}

func GetProdBasisInfo() ProdModel {
	m := GetProdModel()
	if m == nil {
		return ProdModel{}
	}
	return *m
}

type ProdInfoSetRequest struct {
	Name   string `json:"name"`
	Footer string `json:"footer"`
}

func SetProdBasisInfo(req *ProdInfoSetRequest) error {
	m := GetProdModel()
	if m == nil {
		m = NewProdModel()
	}
	m.Name = req.Name
	m.Footer = req.Footer
	return m.Save()
}

func ToggleExperimentalSetting() {
	gNavConf.Experimental = !gNavConf.Experimental
}

func Init() error {
	m := GetProdModel()
	if m != nil {
		return nil
	}

	m = &ProdModel{
		Name:   "xsnos",
		Footer: "xsnos",
	}
	return m.Save()
}

func init() {
	initial.Register("prod-model", initial.DefaultPriority, Init)
}
