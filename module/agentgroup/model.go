package agentgroup

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

// Model  model
type Model struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	isNew       bool
}

func (m *Model) Save() error {
	db := sqldb.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := api.GetColumn(*m, []string{})
		values := api.GetColumnNamed(columns)

		insertClause := fmt.Sprintf("insert into agent_group(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)

		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
		m.isNew = false
	} else {
		columns := api.GetColumn(*m, []string{"id", "created_time"})

		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update agent_group set %s where id=:id", strings.Join(updates, ","))
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.NamedExec(updateClause, m)
		if err != nil {
			logger.Errorf("update model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}
	go SetCache(m)

	return nil
}

func (m *Model) Remove() error {
	db := sqldb.GetDB()

	deleteClause := fmt.Sprintf("delete from agent_group where id=$1")
	_, err := db.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func New() (m *Model) {
	m = new(Model)
	m.ID = uuid.NewString()
	m.isNew = true
	return m
}
