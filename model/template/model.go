package template

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/storage/db"

	msg "github.com/open-cmi/cmmns/msg/template"
)

// Model  model
type Model struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	isNew bool
}

// GetPassword 获取敏感使用
func (m *Model) GetPassword() {
	// 这里获取密码等敏感信息
}

func (m *Model) Save() error {
	dbsql := db.GetDB()

	if m.isNew {
		// 存储到数据库
		id := uuid.New()
		insertClause := fmt.Sprintf("insert into template(id, name) values($1, $2)")

		_, err := dbsql.Exec(insertClause, id.String(), m.Name)
		if err != nil {
			return errors.New("create model failed")
		}
	} else {
		updateClause := fmt.Sprintf("update template set name=$1 where id=$2")
		_, err := dbsql.Exec(updateClause, m.Name, m.ID)
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

func New(reqMsg *msg.CreateMsg) (m *Model) {
	return &Model{
		Name:  reqMsg.Name,
		isNew: true,
	}
}
