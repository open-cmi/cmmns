package template

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/storage/db"
)

// Model  model
type Model struct {
	ID          string `json:"id"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
	UpdatedTime int64  `json:"updated_time" db:"updated_time"`
	Name        string `json:"name"`
	isNew       bool
}

// GetPassword 获取敏感使用
func (m *Model) GetPassword() {
	// 这里获取密码等敏感信息
}

func (m *Model) Save() error {
	dbsql := db.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := []string{"id", "name"}
		var values []string = []string{}
		for index, _ := range columns {
			seq := fmt.Sprintf(`$%d`, index+1)
			values = append(values, seq)
		}

		insertClause := fmt.Sprintf("insert into template(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))
		_, err := dbsql.Exec(insertClause, m.ID, m.Name)
		if err != nil {
			return errors.New("create model failed")
		}

	} else {
		columns := []string{"name"}
		var updates []string = []string{}
		for index, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=%d`, column, index+1))
		}
		updateClause := fmt.Sprintf("update template set %s where id=$%d", strings.Join(updates, ","), len(columns)+1)
		_, err := dbsql.Exec(updateClause, m.Name)
		if err != nil {
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	dbsql := db.GetDB()

	deleteClause := fmt.Sprintf("delete from template where id=$1")
	_, err := dbsql.Exec(deleteClause, m.ID)
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
