package rbac

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

type Role struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
	UpdatedTime int64  `json:"updated_time" db:"updated_time"`
	Permisions  string `json:"permisions" db:"permisions"`
	Description string `json:"description" db:"description"`
	isNew       bool
}

func (r *Role) Save() error {
	db := sqldb.GetDB()

	if r.isNew {
		columns := goparam.GetColumn(Role{}, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into roles(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, r)
		if err != nil {
			return errors.New("create role failed")
		}
	} else {
		columns := goparam.GetColumn(Role{}, []string{})
		values := goparam.GetColumnUpdateNamed(columns)

		insertClause := fmt.Sprintf("update roles set %s where id=:id",
			strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, r)
		if err != nil {
			return errors.New("update user failed")
		}
	}
	return nil
}

func (r *Role) Remove() error {
	deleteClause := "delete from roles where id=$1"
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, r.ID)
	if err != nil {
		return errors.New("remove role failed")
	}
	return err
}

func GetByName(name string) *Role {
	queryClause := `select * from roles where name=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, name)
	if row == nil {
		return nil
	}
	var r Role
	err := row.StructScan(&r)
	if err != nil {
		logger.Errorf("struct scan role %s\n", err.Error())
		return nil
	}
	return &r
}

func GetByID(id string) *Role {
	queryClause := `select * from roles where id=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, id)
	if row == nil {
		return nil
	}
	var r Role
	err := row.StructScan(&r)
	if err != nil {
		logger.Errorf("struct scan role %s\n", err.Error())
		return nil
	}
	return &r
}
