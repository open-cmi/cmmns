package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

type TokenRecord struct {
	Name        string `json:"name" db:"name"`
	Token       string `json:"token" db:"token"`
	ExpireDay   int    `json:"expire_day" db:"expire_day"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
	isNew       bool
}

func (t *TokenRecord) Save() error {
	db := sqldb.GetDB()

	if t.isNew {
		columns := goparam.GetColumn(*t, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into token_record(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, t)
		if err != nil {
			return errors.New("create token record failed")
		}
	} else {
		columns := goparam.GetColumn(*t, []string{})
		values := goparam.GetColumnUpdateNamed(columns)

		insertClause := fmt.Sprintf("update token_record set %s where name=:name",
			strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, t)
		if err != nil {
			return errors.New("update token record failed")
		}
	}
	return nil
}

func GetTokenRecordByToken(token string) *TokenRecord {
	queryClause := `select * from token_record where token=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, token)

	var mdl TokenRecord
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Errorf("token %s not found: %s\n", token, err.Error())
		return nil
	}

	return &mdl
}

func NewTokenRecord() *TokenRecord {
	return &TokenRecord{
		isNew:       true,
		CreatedTime: time.Now().Unix(),
	}
}
