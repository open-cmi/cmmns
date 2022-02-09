package agentgroup

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/storage/db"

	msg "github.com/open-cmi/cmmns/msg/agentgroup"
)

// Model  model
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	isNew       bool
}

func (m *Model) Save() error {
	dbsql := db.GetDB()

	if m.isNew {
		// 存储到数据库
		id := uuid.New()
		insertClause := fmt.Sprintf("insert into agent_group(id, name, description) values($1, $2, $3)")

		_, err := dbsql.Exec(insertClause, id.String(), m.Name, m.Description)
		if err != nil {
			return errors.New("create model failed")
		}
	} else {
		updateClause := fmt.Sprintf("update agent_group set name=$1, description=$2 where id=$3")
		_, err := dbsql.Exec(updateClause, m.Name, m.Description, m.ID)
		if err != nil {
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	dbsql := db.GetDB()

	deleteClause := fmt.Sprintf("delete from agent_group where id=$1")
	_, err := dbsql.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func New(reqMsg *msg.CreateMsg) (m *Model) {
	return &Model{
		Name:        reqMsg.Name,
		Description: reqMsg.Description,
		isNew:       true,
	}
}
