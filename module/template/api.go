package template

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/common/parameter"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

func FilterGet(mo *parameter.Option, fields []string, values []interface{}) *Model {
	var key string = ""
	for index, field := range fields {
		if index != 0 {
			key += "."
		}
		key += fmt.Sprintf("%s.%v", field, values[index])
	}

	m := GetCache(key)
	if m != nil {
		return m
	}

	columns := parameter.GetColumn(Model{}, []string{})

	var whereClause string
	for index, field := range fields {
		if index != 0 {
			whereClause += " and "
		} else {
			whereClause += " where "
		}
		whereClause += fmt.Sprintf(`%s=$%d`, field, index+1)
	}

	queryClause := fmt.Sprintf(`select %s from template %s`, strings.Join(columns, ","), whereClause)
	logger.Debugf(queryClause + "\n")
	db := sqldb.GetConfDB()
	row := db.QueryRowx(queryClause, values...)

	var mdl Model
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return &mdl
}

func Get(mo *parameter.Option, field string, value interface{}) *Model {
	if field == "id" {
		mdl := GetCache(value.(string))
		if mdl != nil {
			return mdl
		}
	}
	return FilterGet(mo, []string{field}, []interface{}{value})
}

// List list
func List(option *parameter.Option) (int, []Model, error) {
	db := sqldb.GetConfDB()

	var results []Model = []Model{}

	countClause := "select count(*) from template"
	whereClause, args := parameter.BuildWhereClause(option)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("count failed: %s\n", err.Error())
		return 0, results, errors.New("get count failed")
	}

	columns := parameter.GetColumn(Model{}, []string{})
	queryClause := fmt.Sprintf(`select %s from template`, strings.Join(columns, ","))
	finalClause := parameter.BuildFinalClause(option)
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
func MultiDelete(mo *parameter.Option, ids []string) error {
	db := sqldb.GetConfDB()

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

func Create(mo *parameter.Option, reqMsg *CreateMsg) (m *Model, err error) {
	// 先检查用户名是否存在
	model := FilterGet(mo, []string{"name"}, []interface{}{reqMsg.Name})
	if model != nil {
		// 用户名已经被占用
		return nil, errors.New("name has been used")
	}
	m = New()
	m.Name = reqMsg.Name
	err = m.Save()

	return m, err
}

func Edit(mo *parameter.Option, id string, reqMsg *EditMsg) error {
	m := FilterGet(mo, []string{"id"}, []interface{}{id})
	if m == nil {
		return errors.New("item not exist")
	}
	m.Name = reqMsg.Name

	err := m.Save()
	return err
}

func Delete(mo *parameter.Option, id string) error {
	m := FilterGet(mo, []string{"id"}, []interface{}{id})
	if m == nil {
		return errors.New("item not exist")
	}
	return m.Remove()
}
