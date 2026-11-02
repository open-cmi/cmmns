package rbac

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

type RoleDetail struct {
	RoleModel
	Permissions []string `json:"permissions"`
}

func GetAllRoleNames() ([]string, error) {
	db := sqldb.GetDB()

	var roles []string = []string{}

	queryClause := `select name from roles`
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return roles, nil
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			logger.Errorf("role scan failed %s\n", err.Error())
			break
		}

		roles = append(roles, name)
	}
	return roles, err
}

func RoleList(option *goparam.Param) (int, []RoleDetail, error) {
	db := sqldb.GetDB()

	var roles []RoleDetail = []RoleDetail{}
	countClause := "select count(*) from roles"

	row := db.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("roles list count failed, %s\n", err.Error())
		return 0, roles, errors.New("list count failed")
	}

	queryClause := `select * from roles`
	finalClause := goparam.BuildFinalClause(option, []string{})
	queryClause += finalClause
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return count, roles, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item RoleModel
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("role struct scan failed %s\n", err.Error())
			break
		}

		perms, _ := GetRolePerms(item.Name)
		roles = append(roles, RoleDetail{
			RoleModel:   item,
			Permissions: perms,
		})
	}
	return count, roles, err
}

type RoleDeleteRequest struct {
	ID string `json:"id"`
}

func DeleteRole(id string) error {
	role := GetByID(id)
	if role == nil {
		return errors.New(i18n.Sprintf("role is not existing"))
	}
	if role.Name == "admin" {
		return errors.New(i18n.Sprintf("admin can not be delete"))
	}
	err := role.Remove()
	return err
}

func CreateRole(name string, description string, perms []string) error {
	if name == "" {
		return errors.New(i18n.Sprintf("role name cannot be empty"))
	}
	r := GetByName(name)
	if r != nil {
		return errors.New(i18n.Sprintf("role already exists"))
	}
	now := time.Now().Unix()
	newRole := RoleModel{
		ID:          uuid.New().String(),
		Name:        name,
		CreatedTime: now,
		UpdatedTime: now,
		Description: description,
		Preset:      false,
		isNew:       true,
	}
	err := newRole.Save()
	if err != nil {
		return err
	}
	if len(perms) > 0 {
		err = UpdateRolePerms(name, perms)
		if err != nil {
			return err
		}
	}
	return nil
}

func EditRole(id string, name string, description string, perms []string) error {
	role := GetByID(id)
	if role == nil {
		return errors.New(i18n.Sprintf("role is not exist"))
	}

	if role.Name == "admin" {
		return errors.New(i18n.Sprintf("admin role cannot be edited"))
	}

	now := time.Now().Unix()
	oldName := role.Name

	if name != "" && name != oldName {
		existing := GetByName(name)
		if existing != nil && existing.ID != role.ID {
			return errors.New(i18n.Sprintf("role already exists"))
		}
		role.Name = name
	}
	role.Description = description
	role.UpdatedTime = now
	role.isNew = false

	err := role.Save()
	if err != nil {
		return err
	}

	if oldName != role.Name {
		db := sqldb.GetDB()
		if db != nil {
			_, _ = db.Exec("DELETE FROM role_permissions WHERE role=$1", oldName)
		}
		InvalidateRolePermCache(oldName)
	}

	err = UpdateRolePerms(role.Name, perms)
	if err != nil {
		return err
	}

	return nil
}

func GetRolePermsByID(id string) ([]string, error) {
	role := GetByID(id)
	if role == nil {
		return nil, errors.New(i18n.Sprintf("role is not exist"))
	}
	return GetRolePerms(role.Name)
}
