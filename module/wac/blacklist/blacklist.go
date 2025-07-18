package blacklist

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

type Blacklist struct {
	Address   string `json:"address" db:"address"`
	Timestamp int64  `json:"timestamp" db:"timestamp"`
	isNew     bool
}

func (u *Blacklist) Save() error {
	db := sqldb.GetDB()

	if u.isNew {
		columns := goparam.GetColumn(Blacklist{}, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into wac_blacklist(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, u)
		if err != nil {
			return errors.New("create blacklist failed")
		}
	} else {
		columns := goparam.GetColumn(Blacklist{}, []string{})
		values := goparam.GetColumnUpdateNamed(columns)

		updateClause := fmt.Sprintf("update wac_blacklist set %s where id=:id",
			strings.Join(values, ","))

		_, err := db.NamedExec(updateClause, u)
		if err != nil {
			logger.Errorf("update blacklist failed: %s\n", err.Error())
			return errors.New("update blacklist failed")
		}
	}
	return nil
}

func (m *Blacklist) Remove() error {
	deleteClause := "delete from wac_blacklist where address=$1"
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, m.Address)
	if err != nil {
		return errors.New("del blacklist failed")
	}
	return err
}

// Get address
func Get(address string) (blacklist *Blacklist) {
	queryClause := `select * from wac_blacklist where address=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, address)

	var mdl Blacklist
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Errorf("blacklist %s not found: %s\n", address, err.Error())
		return nil
	}

	return &mdl
}

func New() *Blacklist {
	return &Blacklist{
		isNew: true,
	}
}

func QueryList(query *goparam.Param) (int, []Blacklist, error) {
	db := sqldb.GetDB()

	var users []Blacklist = []Blacklist{}
	countClause := "select count(*) from wac_blacklist"
	row := db.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("wac_blacklist list count failed, %s\n", err.Error())
		return 0, users, errors.New("list count failed")
	}

	queryClause := `select * from wac_blacklist`
	finalClause := goparam.BuildFinalClause(query)
	queryClause += finalClause
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return count, users, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item Blacklist
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("blacklist struct scan failed %s\n", err.Error())
			break
		}

		users = append(users, item)
	}
	return count, users, err
}

func ListAll() ([]Blacklist, error) {
	db := sqldb.GetDB()

	var users []Blacklist = []Blacklist{}

	queryClause := `select * from wac_blacklist`
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return users, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item Blacklist
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("blacklist struct scan failed %s\n", err.Error())
			break
		}

		users = append(users, item)
	}
	return users, err
}
