package template

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/model"
	msg "github.com/open-cmi/cmmns/msg/template"
	"github.com/open-cmi/cmmns/storage/db"
)

type Option struct {
	model.Option
}

func Get(mo *Option, field string, value string) *Model {
	columns := model.GetColumn(Model{}, []string{})

	queryClause := fmt.Sprintf(`select %s from template where %s=$1`, strings.Join(columns, ","), field)
	sqldb := db.GetDB()
	row := sqldb.QueryRowx(queryClause, value)

	var mdl Model
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil
	}

	return &mdl
}

// List list
func List(option *Option) (int, []Model, error) {
	sqldb := db.GetDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from template")
	whereClause, args := model.BuildWhereClause(&option.Option)
	countClause += whereClause
	row := sqldb.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Logger.Error("count failed: %s\n", err.Error())
		return 0, results, errors.New("get count failed")
	}

	columns := model.GetColumn(Model{}, []string{})
	queryClause := fmt.Sprintf(`select %s from template`, strings.Join(columns, ","))
	finalClause := model.BuildFinalClause(&option.Option)
	queryClause += (whereClause + finalClause)
	rows, err := sqldb.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		logger.Logger.Error(err.Error())
		return count, results, nil
	}

	for rows.Next() {
		var item Model
		err := rows.StructScan(&item)
		if err != nil {
			logger.Logger.Error(err.Error())
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

// List list
func MultiDelete(mo *Option, ids []string) error {
	sqldb := db.GetDB()

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
	_, err := sqldb.Exec(deleteClause, args...)
	if err != nil {
		logger.Logger.Error("delete failed: %s\n", err.Error())
		return errors.New("delete failed")
	}
	return nil
}

func Create(mo *Option, reqMsg *msg.CreateMsg) (m *Model, err error) {
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

func Edit(mo *Option, id string, reqMsg *msg.EditMsg) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	m.Name = reqMsg.Name

	err := m.Save()
	return err
}

func Delete(mo *Option, id string) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	return m.Remove()
}
