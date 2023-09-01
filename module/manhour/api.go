package manhour

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

func Get(mo *api.Option, field string, value string) *Model {
	columns := api.GetColumn(Model{}, []string{})

	queryClause := fmt.Sprintf(`select %s from manhour where %s=$1`, strings.Join(columns, ","), field)
	db := sqldb.GetConfDB()
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
	db := sqldb.GetConfDB()

	var results []Model = []Model{}

	countClause := "select count(*) from manhour"
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
	queryClause := fmt.Sprintf(`select %s from manhour`, strings.Join(columns, ","))
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

	deleteClause := fmt.Sprintf("delete from manhour where id in %s", list)
	_, err := db.Exec(deleteClause, args...)
	if err != nil {
		logger.Errorf("delete failed: %s\n", err.Error())
		return errors.New("delete failed")
	}
	return nil
}

func Create(mo *api.Option, reqMsg *CreateMsg) (m *Model, err error) {
	m = New()
	m.Date = reqMsg.Date
	m.StartTime = reqMsg.StartTime
	m.EndTime = reqMsg.EndTime
	m.Content = reqMsg.Content
	err = m.Save()

	return m, err
}

func Edit(mo *api.Option, id string, reqMsg *EditMsg) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	m.Date = reqMsg.Date
	m.StartTime = reqMsg.StartTime
	m.EndTime = reqMsg.EndTime
	m.Content = reqMsg.Content

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
