package scheduler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/common/parameter"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

func Get(mo *parameter.Option, field string, value string) *Job {
	if field == "id" {
		job := getCache(value)
		if job != nil {
			return job
		}
	}

	columns := parameter.GetColumn(Job{}, []string{})

	queryClause := fmt.Sprintf(`select %s from job where %s=$1`, strings.Join(columns, ","), field)
	db := sqldb.GetConfDB()
	row := db.QueryRowx(queryClause, value)

	var mdl Job
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return &mdl
}

// List list
func List(option *parameter.Option) (int, []Job, error) {
	db := sqldb.GetConfDB()

	var results []Job = []Job{}

	countClause := "select count(*) from job"
	whereClause, args := parameter.BuildWhereClause(option)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("count failed: %s\n", err.Error())
		return 0, results, errors.New("get count failed")
	}

	columns := parameter.GetColumn(Job{}, []string{})
	queryClause := fmt.Sprintf(`select %s from job`, strings.Join(columns, ","))
	finalClause := parameter.BuildFinalClause(option)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		logger.Error(err.Error())
		return count, results, nil
	}

	for rows.Next() {
		var item Job
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

	deleteClause := fmt.Sprintf("delete from job where id in %s", list)
	_, err := db.Exec(deleteClause, args...)
	if err != nil {
		logger.Errorf("delete failed: %s\n", err.Error())
		return errors.New("delete failed")
	}
	return nil
}

func Delete(mo *parameter.Option, id string) error {
	m := Get(mo, "id", id)
	if m == nil {
		return errors.New("item not exist")
	}
	return m.Remove()
}
