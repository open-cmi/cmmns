package secretkey

import (
	"errors"
	"fmt"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

func Get(mo *api.Option, id string) *Model {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select id,name,key_type,key_length,comment from secret_key where id=$1")

	var model Model
	db := sqldb.GetConfDB()
	row := db.QueryRow(queryclause, id)
	err := row.Scan(&model.ID, &model.Name, &model.KeyType, &model.KeyLength, &model.Comment)
	if err == nil {
		// 用户名已经被占用
		return &model
	}
	return nil
}

// List list
func List(mo *api.Option) (int, []Model, error) {
	db := sqldb.GetConfDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from secret_key")
	whereClause, args := api.BuildWhereClause(mo)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	queryClause := fmt.Sprintf(`select id,name,key_type,key_length,comment,public_key from secret_key`)
	queryClause += whereClause
	rows, err := db.Query(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, results, nil
	}

	for rows.Next() {
		var item Model
		err := rows.Scan(&item.ID, &item.Name, &item.KeyType, &item.KeyLength, &item.Comment, &item.PublicKey)
		if err != nil {
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

// List list
func MultiDelete(mo *api.Option, ids []string) error {
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

	deleteClause := fmt.Sprintf("delete from secret_key where id in %s", list)
	_, err := db.Exec(deleteClause, args...)
	if err != nil {
		return errors.New("delete item failed")
	}
	return nil
}

func Create(mo *api.Option, reqMsg *CreateMsg) (m *Model, err error) {
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

func Edit(mo *api.Option, id string, reqMsg *EditMsg) error {
	m := Get(mo, id)
	if m == nil {
		return errors.New("item not exist")
	}

	m.Name = reqMsg.Name
	err := m.Save()
	return err
}

func Delete(mo *api.Option, id string) error {
	m := Get(mo, id)
	if m == nil {
		return errors.New("item not exist")
	}
	return m.Remove()
}
