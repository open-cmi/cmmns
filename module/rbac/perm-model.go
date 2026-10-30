package rbac

import (
	"errors"
	"fmt"
	"time"

	"github.com/open-cmi/gobase/essential/sqldb"
)

type PermType string

const (
	PermView   PermType = "view"
	PermWrite  PermType = "write"
	PermDelete PermType = "delete"
	PermImport PermType = "import"
	PermExport PermType = "export"
)

func PermCode(ident string, perm PermType) string {
	return fmt.Sprintf("%s:%s", ident, perm)
}

type PermModel struct {
	Role        string   `json:"role" db:"role"`
	Perm        PermType `json:"perm" db:"perm"`
	CreatedTime int64    `json:"created_time" db:"created_time"`
	isNew       bool     `json:"-"`
}

func GetPermsFromDB(role string) ([]string, error) {
	db := sqldb.GetDB()
	if db == nil {
		return []string{}, errors.New("db is not available")
	}

	rows, err := db.Queryx(`select perm from role_permissions where role=$1`, role)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	perms := []string{}
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return []string{}, err
		}
		if perm != "" {
			perms = append(perms, perm)
		}
	}
	return perms, nil
}

// UpdateRolePerms 以事务方式更新角色的权限列表
func UpdateRolePerms(role string, perms []string) error {
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

	// 1. 删除该角色现有的所有权限
	_, err = tx.Exec(`DELETE FROM role_permissions WHERE role=$1`, role)
	if err != nil {
		return err
	}

	// 2. 批量插入新权限
	now := time.Now().Unix()
	count := 0
	for _, p := range perms {
		if p == "" {
			continue
		}
		_, err = tx.Exec(`INSERT INTO role_permissions (role, perm, created_time) VALUES ($1, $2, $3)`, role, p, now)
		if err != nil {
			return err
		}
		count++
	}

	err = tx.Commit()
	if err == nil {
		InvalidateRolePermCache(role)
	}
	return err
}
