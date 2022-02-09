package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/model"
	msg "github.com/open-cmi/cmmns/msg/agent"
	"github.com/open-cmi/cmmns/storage/db"
)

type Option struct {
	model.Option
}

func Get(mo *Option, field string, value string) *Model {
	columns := model.GetColumn(Model{}, []string{})

	queryClause := fmt.Sprintf(`select %s from agent where %s=$1`, strings.Join(columns, ","), field)
	dbsql := db.GetDB()
	row := dbsql.QueryRowx(queryClause, value)

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
	dbsql := db.GetDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from agent")
	whereClause, args := model.BuildWhereClause(&option.Option)
	countClause += whereClause
	row := dbsql.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	columns := model.GetColumn(Model{}, []string{})
	queryClause := fmt.Sprintf(`select %s from agent`, strings.Join(columns, ","))
	finalClause := model.BuildFinalClause(&option.Option)
	queryClause += (whereClause + finalClause)
	rows, err := dbsql.Queryx(queryClause, args...)
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

	deleteClause := fmt.Sprintf("delete from agent where id in %s", list)
	_, err := dbsql.Exec(deleteClause, args...)
	if err != nil {
		return errors.New("delete item failed")
	}
	return nil
}

func Create(mo *Option, reqMsg *msg.CreateMsg) (m *Model, err error) {
	// 先检查用户名是否存在
	model := Get(mo, "address", reqMsg.Address)
	if model != nil {
		// 用户名已经被占用
		return nil, errors.New("address has been used")
	}
	m = New()
	m.Group = reqMsg.Group
	m.Address = reqMsg.Address
	m.Port = reqMsg.Port
	m.ConnType = reqMsg.ConnType
	m.UserName = reqMsg.UserName
	m.Passwd = reqMsg.Passwd
	m.SecretKey = reqMsg.SecretKey
	m.Description = reqMsg.Description
	err = m.Save()

	return m, err
}

func Edit(mo *Option, id string, reqMsg *msg.EditMsg) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	m.Address = reqMsg.Address
	m.ConnType = reqMsg.ConnType
	m.Description = reqMsg.Description
	m.Group = reqMsg.Group
	m.Passwd = reqMsg.Passwd
	m.Port = reqMsg.Port
	m.SecretKey = reqMsg.SecretKey
	m.UserName = reqMsg.UserName

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
