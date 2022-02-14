package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/storage/sqldb"
	"github.com/open-cmi/cmmns/module/template/msg"
)

func Get(mo *api.Option, field string, value string) *Model {
	columns := api.GetColumn(Model{}, []string{})

	queryClause := fmt.Sprintf(`select %s from template where %s=$1`, strings.Join(columns, ","), field)
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, value)

	var mdl Model
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return &mdl
}

// List list
func List(option *api.Option) (int, []Model, error) {
	db := sqldb.GetDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from template")
	whereClause, args := api.BuildWhereClause(option)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("count failed: %s\n", err.Error())
		return 0, results, errors.New("get count failed")
	}

	columns := api.GetColumn(Model{}, []string{})
	queryClause := fmt.Sprintf(`select %s from template`, strings.Join(columns, ","))
	finalClause := api.BuildFinalClause(option)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		logger.Error(err.Error())
		return count, results, nil
	}

	for rows.Next() {
		var item Model
		err := rows.StructScan(&item)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

// List list
func MultiDelete(mo *api.Option, ids []string) error {
	db := sqldb.GetDB()

	if len(ids) == 0 {
		return errors.New("no items deleted")
	}

	var list = " ("
	for index, _ := range ids {
		if index != 0 {
			list += ","
		}
		list += fmt.Sprintf("$%d", index+1)
	}
	list += ")"

	var args []interface{} = []interface{}{}
	for _, item := range ids {
		args = append(args, item)
	}

	deleteClause := fmt.Sprintf("delete from template where id in %s", list)
	_, err := db.Exec(deleteClause, args...)
	if err != nil {
		logger.Errorf("delete failed: %s\n", err.Error())
		return errors.New("delete failed")
	}
	return nil
}

func Create(mo *api.Option, reqMsg *msg.CreateMsg) (m *Model, err error) {
	// 先检查用户名是否存在
	model := Get(mo, "name", reqMsg.Name)
	if model != nil {
		// 用户名已经被占用
		return nil, errors.New("name has been used")
	}
	m = New()
	m.Name = reqMsg.Name
	err = m.Save()

	return m, err
}

func Edit(mo *api.Option, id string, reqMsg *msg.EditMsg) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	m.Name = reqMsg.Name

	err := m.Save()
	return err
}

func Delete(mo *api.Option, id string) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	return m.Remove()
}
