package email

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

type EmailModel struct {
	ID       string `json:"id" db:"id"`
	Server   string `json:"server" db:"server"`
	Port     int    `json:"port" db:"port"`
	Sender   string `json:"sender" db:"sender"`
	Password string `json:"password" db:"password"`
	UseTLS   bool   `json:"use_tls" db:"use_tls"`
	IsNew    bool
}

func (m *EmailModel) Save() error {
	db := sqldb.GetConfDB()

	if m.IsNew {
		m.ID = uuid.NewString()
		// 存储到数据库
		columns := goparam.GetColumn(*m, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into sender_email(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s", insertClause)

		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
	} else {
		columns := goparam.GetColumn(*m, []string{"id"})

		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update sender_email set %s where id=:id", strings.Join(updates, ","))
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.NamedExec(updateClause, m)
		if err != nil {
			logger.Errorf("update sender_email model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}
	return nil
}

func New() *EmailModel {
	return &EmailModel{
		IsNew: true,
	}
}

func Get() *EmailModel {
	sqlClause := fmt.Sprintf("select * from sender_email")

	clause := sqlClause

	db := sqldb.GetConfDB()
	row := db.QueryRowx(clause)
	if row == nil {
		return nil
	}
	var m EmailModel
	err := row.StructScan(&m)
	if err != nil {
		return nil
	}
	return &m
}
