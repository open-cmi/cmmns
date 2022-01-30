package template

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/model"
	"github.com/open-cmi/cmmns/storage/db"
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
	sqldb := db.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := model.GetColumn(*m, []string{})
		values := model.GetColumnNamed(columns)

		insertClause := fmt.Sprintf("insert into template(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Logger.Info("start to exec sql clause: %s", insertClause)

		_, err := sqldb.NamedExec(insertClause, m)
		if err != nil {
			logger.Logger.Error("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
	} else {
		columns := model.GetColumn(*m, []string{"id", "created_time"})

		m.UpdatedTime = time.Now().Unix()
		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update template set %s where id=:id", strings.Join(updates, ","))
		logger.Logger.Debug("start to exec sql clause: %s", updateClause)
		_, err := sqldb.NamedExec(updateClause, m)
		if err != nil {
			logger.Logger.Error("update template model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	sqldb := db.GetDB()

	deleteClause := fmt.Sprintf("delete from template where id=$1")
	_, err := sqldb.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func New() (m *Model) {
	now := time.Now().Unix()
	return &Model{
		ID:          uuid.NewString(),
		CreatedTime: now,
		UpdatedTime: now,
		isNew:       true,
	}
}
