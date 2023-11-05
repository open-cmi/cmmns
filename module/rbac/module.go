package rbac

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

type Module struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	isNew       bool
}

func (m *Module) Save() error {
	db := sqldb.GetConfDB()

	if m.isNew {
		columns := goparam.GetColumn(Module{}, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into modules(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			return errors.New("create role failed")
		}
	} else {
		columns := goparam.GetColumn(Module{}, []string{})
		values := goparam.GetColumnUpdateNamed(columns)

		insertClause := fmt.Sprintf("update modules set %s where name=:name",
			strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			return errors.New("update user failed")
		}
	}
	return nil
}

func NewModule(name string, desc string) *Module {
	return &Module{
		Name:        name,
		Description: desc,
		isNew:       true,
	}
}

func GetModule(name string) *Module {
	queryClause := `select * from modules where name=$1`
	db := sqldb.GetConfDB()

	row := db.QueryRowx(queryClause, name)
	if row == nil {
		return nil
	}
	var m Module
	err := row.StructScan(&m)
	if err != nil {
		logger.Errorf("struct scan module failed: %s\n", err.Error())
		return nil
	}
	return &m
}
