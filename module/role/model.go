package role

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
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
	db := sqldb.GetConfDB()

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

func (r *Role) HasReadPermision(m string) bool {
	return r.hasPermision(m, "read")
}

func (r *Role) HasWritePermision(m string) bool {
	return r.hasPermision(m, "write")
}

func (r *Role) hasPermision(m string, perm string) bool {
	if r.Permisions == "*" {
		return true
	}

	var hasPerm bool
	mperms := strings.Split(r.Permisions, ";")
	for _, mperm := range mperms {
		s := strings.Split(mperm, ":")
		mod := s[0]
		permision := s[1]
		if mod != m {
			continue
		}
		perms := strings.Split(permision, ",")
		for _, p := range perms {
			if perm == p {
				hasPerm = true
				break
			}
		}
	}
	return hasPerm
}

func Get(name string) *Role {
	queryClause := `select * from roles where name=$1`
	db := sqldb.GetConfDB()
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
