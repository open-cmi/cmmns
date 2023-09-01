package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

func FilterGet(mo *api.Option, fields []string, values []interface{}) *Model {
	columns := api.GetColumn(Model{}, []string{})

	var whereClause string
	for index, field := range fields {
		if index != 0 {
			whereClause += " and "
		} else {
			whereClause += " where "
		}
		whereClause += fmt.Sprintf(`%s=$%d`, field, index+1)
	}

	queryClause := fmt.Sprintf(`select %s from agent %s`, strings.Join(columns, ","), whereClause)
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

func Get(field string, value interface{}) *Model {
	if field == "dev_id" {
		mdl := GetCache(value.(string))
		if mdl != nil {
			return mdl
		}
	}
	return FilterGet(nil, []string{field}, []interface{}{value})
}

// List list
func List(option *api.Option) (int, []Model, error) {
	db := sqldb.GetConfDB()

	var results []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from agent")
	whereClause, args := api.BuildWhereClause(option)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	columns := api.GetColumn(Model{}, []string{})
	queryClause := fmt.Sprintf(`select %s from agent`, strings.Join(columns, ","))
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

	deleteClause := fmt.Sprintf("delete from agent where id in %s", list)
	_, err := db.Exec(deleteClause, args...)
	if err != nil {
		return errors.New("delete item failed")
	}
	return nil
}

func Create(mo *api.Option, reqMsg *CreateMsg) (m *Model, err error) {
	// 先检查用户名是否存在
	api := FilterGet(mo, []string{"address"}, []interface{}{reqMsg.Address})
	if api != nil {
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

func Edit(mo *api.Option, id string, reqMsg *EditMsg) error {
	m := FilterGet(mo, []string{"id"}, []interface{}{id})
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

func Delete(mo *api.Option, id string) error {
	m := FilterGet(mo, []string{"id"}, []interface{}{id})
	if m == nil {
		return errors.New("item not exist")
	}

	DeleteCache(m.DevID)
	return m.Remove()
}
