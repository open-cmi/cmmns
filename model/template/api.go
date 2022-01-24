package template

import (
	"errors"
	"fmt"

	"github.com/open-cmi/cmmns/msg/request"
	msg "github.com/open-cmi/cmmns/msg/template"
	"github.com/open-cmi/cmmns/storage/db"
	"github.com/open-cmi/cmmns/utils"
)

type ModelOption struct {
	UserID string
}

func Get(mo *ModelOption, field string, value string) *Model {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select * from template where %s=$1", field)

	var model Model
	sqldb := db.GetDB()
	row := sqldb.QueryRow(queryclause, value)
	err := row.Scan(&model.ID, &model.Name)
	if err == nil {
		// 用户名已经被占用
		return &model
	}
	return nil
}

// List list
func List(mo *ModelOption, p *request.RequestQuery) (int, []Model, error) {
	dbsql := db.GetDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from template")
	whereClause, args := utils.BuildWhereClause(p)
	countClause += whereClause
	row := dbsql.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	queryClause := fmt.Sprintf(`select id,name from template`)
	finalClause := utils.BuildFinalClause(p)
	queryClause += (whereClause + finalClause)
	rows, err := dbsql.Query(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, results, nil
	}

	for rows.Next() {
		var item Model
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

// List list
func MultiDelete(mo *ModelOption, ids []string) error {
	dbsql := db.GetDB()

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
	_, err := dbsql.Exec(deleteClause, args...)
	if err != nil {
		return errors.New("delete item failed")
	}
	return nil
}

func Create(mo *ModelOption, reqMsg *msg.CreateMsg) (m *Model, err error) {
	// 先检查用户名是否存在
	model := Get(mo, "name", reqMsg.Name)
	if model != nil {
		// 用户名已经被占用
		return nil, errors.New("name has been used")
	}

	m = New(reqMsg)
	err = m.Save()

	return m, err
}

func Edit(mo *ModelOption, id string, reqMsg *msg.EditMsg) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	m.Name = reqMsg.Name
	err := m.Save()
	return err
}

func Delete(mo *ModelOption, id string) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	return m.Remove()
}
