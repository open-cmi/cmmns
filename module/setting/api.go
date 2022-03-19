package setting

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

func FilterGet(mo *api.Option, fields []string, values []interface{}) *Model {
	var key string = ""
	for index, field := range fields {
		if index != 0 {
			key += "."
		}
		key += fmt.Sprintf("%s.%v", field, values[index])
	}

	mdl := GetCache(key)
	if mdl != nil {
		return mdl
	}

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

	queryClause := fmt.Sprintf(`select %s from setting %s`, strings.Join(columns, ","), whereClause)
	logger.Debugf(queryClause + "\n")
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, values...)

	var model Model
	err := row.StructScan(&model)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	go SetCache(key, &model)
	return &model
}

func Get(mo *api.Option, field string, value interface{}) *Model {
	return FilterGet(mo, []string{field}, []interface{}{value})
}

// List list
func List(option *api.Option) (int, []Model, error) {
	db := sqldb.GetDB()

	var results []Model = []Model{}

	countClause := "select count(*) from setting"
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
	queryClause := fmt.Sprintf(`select %s from setting`, strings.Join(columns, ","))
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

func Edit(mo *api.Option, id string, reqMsg *EditMsg) error {
	m := FilterGet(mo, []string{"id"}, []interface{}{id})
	if m == nil {
		return errors.New("item not exist")
	}

	m.Scope = reqMsg.Scope
	m.Belong = reqMsg.Belong
	m.Key = reqMsg.Key
	m.Value = reqMsg.Value
	m.CfgSeq++

	err := m.Save()
	return err
}
