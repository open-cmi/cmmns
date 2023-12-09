package rbac

import (
	"errors"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

func RoleNameList() (int, []string, error) {
	db := sqldb.GetConfDB()

	var roles []string = []string{}
	countClause := "select count(*) from roles"

	row := db.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("roles list count failed, %s\n", err.Error())
		return 0, roles, errors.New("list count failed")
	}

	queryClause := `select name from roles`
	rows, err := db.Queryx(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return count, roles, nil
	}

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			logger.Errorf("role scan failed %s\n", err.Error())
			break
		}

		roles = append(roles, name)
	}
	return count, roles, err
}

func RoleList(option *goparam.Option) (int, []Role, error) {
	db := sqldb.GetConfDB()

	var roles []Role = []Role{}
	countClause := "select count(*) from roles"

	whereClause, args := goparam.BuildWhereClause(option)

	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("roles list count failed, %s\n", err.Error())
		return 0, roles, errors.New("list count failed")
	}

	queryClause := `select * from roles`
	finalClause := goparam.BuildFinalClause(option)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, roles, nil
	}

	for rows.Next() {
		var item Role
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("role struct scan failed %s\n", err.Error())
			break
		}

		roles = append(roles, item)
	}
	return count, roles, err
}

func ModuleList(option *goparam.Option) (int, []Module, error) {
	db := sqldb.GetConfDB()

	var modules []Module = []Module{}
	countClause := "select count(*) from modules"

	whereClause, args := goparam.BuildWhereClause(option)

	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("modules list count failed, %s\n", err.Error())
		return 0, modules, errors.New("list count failed")
	}

	queryClause := `select * from modules`
	finalClause := goparam.BuildFinalClause(option)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, modules, nil
	}

	for rows.Next() {
		var item Module
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("module struct scan failed %s\n", err.Error())
			break
		}

		modules = append(modules, item)
	}
	return count, modules, err
}
