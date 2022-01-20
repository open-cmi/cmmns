package secretkey

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/storage/db"

	msg "github.com/open-cmi/cmmns/msg/secretkey"
)

// Model  model
type Model struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	KeyType      string `json:"key_type"`
	KeyLength    int    `json:"key_length"`
	Comment      string `json:"comment"`
	PassPrase    string `json:"passprase"`
	Confirmation string `json:"confirmation"`
	isNew        bool
}

func (m *Model) Save() error {
	dbsql := db.GetDB()

	if m.isNew {
		// 存储到数据库
		id := uuid.New()
		insertClause := fmt.Sprintf(`insert into secret_key(id, name, key_type, key_length, comment, passprase, confirmation) 
		values($1, $2, $3, $4, $5, $6, $7)`)

		_, err := dbsql.Exec(insertClause, id.String(), m.Name, m.KeyType, m.KeyLength, m.Comment, m.PassPrase, m.Confirmation)
		if err != nil {
			return errors.New("create model failed")
		}
	} else {
		updateClause := fmt.Sprintf("update secret_key set name=$1,key_type=$2,key_length=$3,comment=$4,passprase=$5,confirmation=$6 where id=$7")
		_, err := dbsql.Exec(updateClause, m.Name, m.ID)
		if err != nil {
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	dbsql := db.GetDB()

	deleteClause := fmt.Sprintf("delete from secret_key where id=$1")
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
