package rbac

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

type RoleModel struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	CreatedTime int64  `json:"created_time" db:"created_time"`
	UpdatedTime int64  `json:"updated_time" db:"updated_time"`
	Description string `json:"description" db:"description"`
	Preset      bool   `json:"preset" db:"preset"`
	isNew       bool
}

func (r *RoleModel) Save() error {
	db := sqldb.GetDB()

	if r.isNew {
		columns := goparam.GetColumn(RoleModel{}, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into roles(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, r)
		if err != nil {
			return errors.New("create role failed")
		}
	} else {
		columns := goparam.GetColumn(RoleModel{}, []string{})
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

func (r *RoleModel) Remove() error {
	db := sqldb.GetDB()
	if db == nil {
		return errors.New("db is not available")
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. 检查是否有用户仍属于该角色，存在关联用户时禁止删除
	var userCount int
	err = tx.QueryRow("SELECT count(*) FROM users WHERE role=$1", r.Name).Scan(&userCount)
	if err != nil {
		logger.Errorf("count users with role %s failed: %s\n", r.Name, err.Error())
		return errors.New("check role users failed")
	}
	if userCount > 0 {
		return errors.New(i18n.Sprintf("role '%s' is in use by %d user(s) and cannot be deleted", r.Name, userCount))
	}

	// 2. 级联删除该角色关联的权限记录
	_, err = tx.Exec("DELETE FROM role_permissions WHERE role=$1", r.Name)
	if err != nil {
		logger.Errorf("delete role_permissions with role %s failed: %s\n", r.Name, err.Error())
		return errors.New("remove role permissions failed")
	}

	// 3. 删除角色记录
	_, err = tx.Exec("DELETE FROM roles WHERE id=$1", r.ID)
	if err != nil {
		logger.Errorf("delete role %s failed: %s\n", r.ID, err.Error())
		return errors.New("remove role failed")
	}

	err = tx.Commit()
	if err == nil {
		InvalidateRolePermCache(r.Name)
	}
	return err
}

func GetByName(name string) *RoleModel {
	queryClause := `select * from roles where name=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, name)
	if row == nil {
		return nil
	}
	var r RoleModel
	err := row.StructScan(&r)
	if err != nil {
		logger.Errorf("struct scan role %s\n", err.Error())
		return nil
	}
	return &r
}

func GetByID(id string) *RoleModel {
	queryClause := `select * from roles where id=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, id)
	if row == nil {
		return nil
	}
	var r RoleModel
	err := row.StructScan(&r)
	if err != nil {
		logger.Errorf("struct scan role %s\n", err.Error())
		return nil
	}
	return &r
}
