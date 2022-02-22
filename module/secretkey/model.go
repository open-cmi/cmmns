package secretkey

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/open-cmi/cmmns/essential/sqldb"
)

// Model  model
type Model struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	KeyType      string `json:"key_type"`
	KeyLength    int    `json:"key_length"`
	Comment      string `json:"comment"`
	PassPhrase   string `json:"passphrase"`
	Confirmation string `json:"confirmation"`
	PrivateKey   string `json:"private_key"`
	PublicKey    string `json:"public_key"`
	isNew        bool
}

func (m *Model) Save() error {
	db := sqldb.GetDB()

	if m.isNew {
		// 存储到数据库
		id := uuid.New()
		insertClause := fmt.Sprintf(`insert into 
			secret_key(id, name, key_type, key_length, comment, passphrase, confirmation, private_key, public_key) 
			values($1, $2, $3, $4, $5, $6, $7, $8, $9)`)
		_, err := db.Exec(insertClause, id.String(), m.Name, m.KeyType,
			m.KeyLength, m.Comment, m.PassPhrase, m.Confirmation, m.PrivateKey, m.PublicKey)
		if err != nil {
			return errors.New("create model failed")
		}
	} else {
		updateClause := fmt.Sprintf(`update secret_key set name=$1 where id=$2`)
		_, err := db.Exec(updateClause, m.Name, m.ID)
		if err != nil {
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	db := sqldb.GetDB()

	deleteClause := fmt.Sprintf("delete from secret_key where id=$1")
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
		Name:         reqMsg.Name,
		KeyType:      reqMsg.KeyType,
		KeyLength:    reqMsg.KeyLength,
		Comment:      reqMsg.Comment,
		PassPhrase:   reqMsg.PassPhrase,
		Confirmation: reqMsg.Confirmation,
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
		isNew:        true,
	}
}
