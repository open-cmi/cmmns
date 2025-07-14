package token

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/essential/webserver/middleware"
	"github.com/open-cmi/gobase/pkg/goparam"
)

// List list
func QueryTokenList(option *goparam.Param) (int, []TokenRecord, error) {
	db := sqldb.GetDB()

	var results []TokenRecord = []TokenRecord{}

	countClause := "select count(*) from token_record"
	row := db.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("count failed: %s\n", err.Error())
		return 0, results, errors.New("get count failed")
	}

	columns := goparam.GetColumn(TokenRecord{}, []string{})
	queryClause := fmt.Sprintf(`select %s from token_record`, strings.Join(columns, ","))
	finalClause := goparam.BuildFinalClause(option)
	queryClause += finalClause
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		logger.Error(err.Error())
		return count, results, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item TokenRecord
		err := rows.StructScan(&item)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

type DeleteTokenRequest struct {
	Name string `json:"name"`
}

func DeleteAuthToken(name string) error {
	t := GetTokenRecord(name)
	if t == nil {
		return errors.New("token is not existing")
	}
	err := t.Remove()
	if err != nil {
		return err
	}
	return nil
}

type CreateTokenRequest struct {
	Name      string `json:"name"`
	ExpireDay int    `json:"expire_day"`
}

func CreateToken(name string, username string, id string, email string, role int, status int, expireDay int) error {
	token, err := middleware.GenerateAuthToken(username, id, email, role, status, expireDay)
	if err != nil {
		return err
	}

	t := NewTokenRecord()
	t.ExpireDay = expireDay
	t.Token = token
	t.Name = name
	err = t.Save()
	if err != nil {
		return err
	}
	return nil
}
