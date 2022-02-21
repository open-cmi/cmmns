package model

import (
	"errors"
	"fmt"

	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/pubsub"
	"github.com/open-cmi/cmmns/essential/storage/sqldb"
	"github.com/open-cmi/cmmns/module/agentgroup/msg"
)

func Get(mo *api.Option, id string) *Model {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select * from agent_group where id=$1")

	var model Model
	db := sqldb.GetDB()
	row := db.QueryRow(queryclause, id)
	err := row.Scan(&model.ID, &model.Name)
	if err == nil {
		// 用户名已经被占用
		return &model
	}
	return nil
}

// List list
func List(mo *api.Option) (int, []Model, error) {
	db := sqldb.GetDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from agent_group")
	whereClause, args := api.BuildWhereClause(mo)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	queryClause := fmt.Sprintf(`select id,name,description from agent_group`)
	queryClause += whereClause
	rows, err := db.Query(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, results, nil
	}

	for rows.Next() {
		var item Model
		err := rows.Scan(&item.ID, &item.Name, &item.Description)
		if err != nil {
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

	deleteClause := fmt.Sprintf("delete from agent_group where id in %s", list)
	_, err := db.Exec(deleteClause, args...)
	if err != nil {
		return errors.New("delete item failed")
	}
	return nil
}

func Create(mo *api.Option, reqMsg *msg.CreateMsg) (m *Model, err error) {
	// 先检查用户名是否存在
	model := Get(mo, reqMsg.Name)
	if model != nil {
		// 用户名已经被占用
		return nil, errors.New("name has been used")
	}

	m = New(reqMsg)
	err = m.Save()

	return m, err
}

func Edit(mo *api.Option, name string, reqMsg *msg.EditMsg) error {
	m := Get(mo, name)
	if m == nil {
		return errors.New("item not exist")
	}

	m.Name = reqMsg.Name
	err := m.Save()
	return err
}

func Delete(mo *api.Option, name string) error {
	m := Get(mo, name)
	if m == nil {
		return errors.New("item not exist")
	}
	return m.Remove()
}

func init() {
	pubsub.Subscribe(def.EventUserCreate, func(username string) {
		fmt.Println("user create:", username)
	})
}
