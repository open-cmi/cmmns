package secretkey

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

// Model  model
type Model struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	KeyType      string `json:"key_type" db:"key_type"`
	KeyLength    int    `json:"key_length" db:"key_length"`
	Comment      string `json:"comment" db:"comment"`
	PassPhrase   string `json:"passphrase" db:"passphrase"`
	Confirmation string `json:"confirmation" db:"confirmation"`
	PrivateKey   string `json:"-" db:"private_key"`
	PublicKey    string `json:"public_key" db:"public_key"`
	CreatedTime  int64  `json:"created_time" db:"created_time"`
	isNew        bool
}

func (m *Model) Save() error {
	db := sqldb.GetConfDB()

	if m.isNew {
		// 存储到数据库
		columns := goparam.GetColumn(*m, []string{})
		insertedColumn := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf(`insert into secret_key(%s) values(%s)`,
			strings.Join(columns, ","), strings.Join(insertedColumn, ","))
		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			return errors.New("create model failed")
		}
	} else {
		columns := goparam.GetColumn(*m, []string{})
		updateColumns := goparam.GetColumnUpdateNamed(columns)
		updateClause := fmt.Sprintf(`update secret_key set %s where id=:id`, strings.Join(updateColumns, ","))
		_, err := db.NamedExec(updateClause, m)
		if err != nil {
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	db := sqldb.GetConfDB()

	deleteClause := "delete from secret_key where id=$1"
	_, err := db.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func New(reqMsg *CreateMsg) (m *Model) {
	privateKey, publicKey, _ := GenerateSecretKey(reqMsg.Name, reqMsg.KeyType,
		reqMsg.KeyLength, reqMsg.Comment, reqMsg.PassPhrase)

	return &Model{
		ID:           uuid.NewString(),
		Name:         reqMsg.Name,
		KeyType:      reqMsg.KeyType,
		KeyLength:    reqMsg.KeyLength,
		Comment:      reqMsg.Comment,
		PassPhrase:   reqMsg.PassPhrase,
		Confirmation: reqMsg.Confirmation,
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
		CreatedTime:  time.Now().Unix(),
		isNew:        true,
	}
}
