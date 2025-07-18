package licmng

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

type LicenseModel struct {
	ID          string `json:"id" db:"id"`
	Customer    string `json:"customer" db:"customer"`
	Prod        string `json:"prod" db:"prod"`
	Version     string `json:"version" db:"version"`
	Modules     string `json:"modules" db:"modules"`
	ExpireTime  int64  `json:"expire_time" db:"expire_time"`
	MCode       string `json:"mcode" db:"mcode"`
	Serial      string `json:"serial"`
	Model       string `json:"model" db:"model"`
	Username    string `json:"username" db:"username"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
	UpdatedTime int64  `json:"updated_time" db:"updated_time"`
	isNew       bool
}

func (m *LicenseModel) Save() error {
	db := sqldb.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := goparam.GetColumn(*m, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into license(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)

		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
		m.isNew = false
	} else {
		columns := goparam.GetColumn(*m, []string{"id", "created_time"})

		m.UpdatedTime = time.Now().Unix()
		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update license set %s where id=:id", strings.Join(updates, ","))
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.NamedExec(updateClause, m)
		if err != nil {
			logger.Errorf("update license model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *LicenseModel) Remove() error {
	db := sqldb.GetDB()

	deleteClause := "delete from license where id=$1"
	_, err := db.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func New() *LicenseModel {
	return &LicenseModel{
		ID:          uuid.New().String(),
		CreatedTime: time.Now().Unix(),
		UpdatedTime: time.Now().Unix(),
		isNew:       true,
	}
}

func Get(id string) *LicenseModel {
	queryClause := `select * from license where id=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, id)

	var mdl LicenseModel
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Errorf("license %s not found: %s\n", id, err.Error())
		return nil
	}

	return &mdl
}
