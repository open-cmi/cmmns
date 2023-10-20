package template

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/common/parameter"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

// Model  model
type Model struct {
	ID          string `json:"id" db:"id"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
	UpdatedTime int64  `json:"updated_time" db:"updated_time"`
	Name        string `json:"name" db:"name"`
	isNew       bool
}

func (m *Model) Save() error {
	db := sqldb.GetConfDB()

	if m.isNew {
		// 存储到数据库
		columns := parameter.GetColumn(*m, []string{})
		values := parameter.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into template(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)

		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
		m.isNew = false
	} else {
		columns := parameter.GetColumn(*m, []string{"id", "created_time"})

		m.UpdatedTime = time.Now().Unix()
		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update template set %s where id=:id", strings.Join(updates, ","))
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.NamedExec(updateClause, m)
		if err != nil {
			logger.Errorf("update template model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}
	go SetCache(m)

	return nil
}

func (m *Model) Remove() error {
	db := sqldb.GetConfDB()

	deleteClause := fmt.Sprintf("delete from template where id=$1")
	_, err := db.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func New() (m *Model) {
	now := time.Now().Unix()
	m = new(Model)
	m.ID = uuid.NewString()
	m.CreatedTime = now
	m.UpdatedTime = now
	m.isNew = true
	return m
}
