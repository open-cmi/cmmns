package whitelist

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

type Whitelist struct {
	Address   string `json:"address" db:"address"`
	Timestamp int64  `json:"timestamp" db:"timestamp"`
	isNew     bool
}

func (u *Whitelist) Save() error {
	db := sqldb.GetDB()

	if u.isNew {
		columns := goparam.GetColumn(Whitelist{}, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into wac_whitelist(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, u)
		if err != nil {
			return errors.New("create whitelist failed")
		}
	} else {
		columns := goparam.GetColumn(Whitelist{}, []string{})
		values := goparam.GetColumnUpdateNamed(columns)

		updateClause := fmt.Sprintf("update wac_whitelist set %s where id=:id",
			strings.Join(values, ","))

		_, err := db.NamedExec(updateClause, u)
		if err != nil {
			logger.Errorf("update whitelist failed: %s\n", err.Error())
			return errors.New("update whitelist failed")
		}
	}
	return nil
}

func (m *Whitelist) Remove() error {
	deleteClause := "delete from wac_whitelist where address=$1"
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, m.Address)
	if err != nil {
		return errors.New("del whitelist failed")
	}
	return err
}

// Get address
func Get(address string) (whitelist *Whitelist) {
	queryClause := `select * from wac_whitelist where address=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, address)

	var mdl Whitelist
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Errorf("whitelist %s not found: %s\n", address, err.Error())
		return nil
	}

	return &mdl
}

func New() *Whitelist {
	return &Whitelist{
		isNew: true,
	}
}

func List(query *goparam.Param) (int, []Whitelist, error) {
	db := sqldb.GetDB()

	var users []Whitelist = []Whitelist{}
	countClause := "select count(*) from wac_whitelist"

	whereClause, args := goparam.BuildWhereClause(query)

	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("wac_whitelist list count failed, %s\n", err.Error())
		return 0, users, errors.New("list count failed")
	}

	queryClause := `select * from wac_whitelist`
	finalClause := goparam.BuildFinalClause(query)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, users, nil
	}

	for rows.Next() {
		var item Whitelist
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("whitelist struct scan failed %s\n", err.Error())
			break
		}

		users = append(users, item)
	}
	return count, users, err
}
