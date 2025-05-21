package time

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

func CreateAbsoluteTimeObject(req *AbsoluteTimeObject) error {
	m := GetTimeObject(req.Name)
	if m != nil {
		return fmt.Errorf("time object with %s is exist", req.Name)
	}

	if req.TimestampStart >= req.TimestampEnd {
		return fmt.Errorf("start time must smaller than end time")
	}

	v, _ := json.Marshal(req)

	obj := CreateNewTimeObject(req.Name, req.Description, TimeTypeAbsolute, string(v))

	return obj.Save()
}

func CreatePeriodTimeObject(req *PeriodTimeObject) error {
	m := GetTimeObject(req.Name)
	if m != nil {
		return fmt.Errorf("time object with %s is exist", req.Name)
	}

	if req.TimeRangeEnable {
		if req.TimeRangeStart > req.TimeRangeEnd {
			return fmt.Errorf("range start time must smaller than range end time")
		}
	}

	for _, period := range req.WeekPeriods {
		if period.PeriodTimeStart > period.PeriodTimeEnd {
			return fmt.Errorf("period start time must smaller than period end time")
		}
	}
	v, _ := json.Marshal(req)
	m = CreateNewTimeObject(req.Name, req.Description, TimeTypePeriod, string(v))

	return m.Save()
}

func QueryAbsoluteTimeObjectList(param *goparam.Param) (int, []AbsoluteTimeObject, error) {
	db := sqldb.GetDB()

	var results []AbsoluteTimeObject = []AbsoluteTimeObject{}

	countClause := "select count(*) from object_time"
	if param.WhereClause != "" {
		param.WhereArgs = append(param.WhereArgs, TimeTypeAbsolute)
		param.WhereClause += fmt.Sprintf(" and time_type=$%d", len(param.WhereArgs))
	} else {
		param.WhereArgs = append(param.WhereArgs, TimeTypeAbsolute)
		param.WhereClause += fmt.Sprintf(" where time_type=$%d", len(param.WhereArgs))
	}
	countClause += param.WhereClause
	row := db.QueryRow(countClause, param.WhereArgs...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("count failed: %s\n", err.Error())
		return 0, results, errors.New("get count failed")
	}

	queryClause := `select * from object_time`
	finalClause := goparam.BuildFinalClause(param)
	queryClause += (param.WhereClause + finalClause)
	rows, err := db.Queryx(queryClause, param.WhereArgs...)
	if err != nil {
		logger.Errorf("time object queryx failed: %s\n", err.Error())
		return 0, results, errors.New("no rows found")
	}
	defer rows.Close()
	for rows.Next() {
		var mdl TimeObject
		err := rows.StructScan(&mdl)
		if err != nil {
			logger.Error(err.Error())
			break
		}
		var obj AbsoluteTimeObject
		err = json.Unmarshal([]byte(mdl.Value), &obj)
		if err != nil {
			logger.Errorf("query abs time list unmarshal failed: %s\n", err.Error())
			break
		}
		obj.Active = obj.IsActive()
		results = append(results, obj)
	}
	return count, results, err
}

func QueryPeriodTimeObjectList(param *goparam.Param) (int, []PeriodTimeObject, error) {
	db := sqldb.GetDB()

	var results []PeriodTimeObject = []PeriodTimeObject{}

	countClause := "select count(*) from object_time"
	if param.WhereClause != "" {
		param.WhereArgs = append(param.WhereArgs, TimeTypePeriod)
		param.WhereClause += fmt.Sprintf(" and time_type=$%d", len(param.WhereArgs))
	} else {
		param.WhereArgs = append(param.WhereArgs, TimeTypePeriod)
		param.WhereClause += fmt.Sprintf(" where time_type=$%d", len(param.WhereArgs))
	}
	countClause += param.WhereClause
	row := db.QueryRow(countClause, param.WhereArgs...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("count failed: %s\n", err.Error())
		return 0, results, errors.New("get count failed")
	}

	queryClause := `select * from object_time`
	finalClause := goparam.BuildFinalClause(param)
	queryClause += (param.WhereClause + finalClause)
	rows, err := db.Queryx(queryClause, param.WhereArgs...)
	if err != nil {
		logger.Errorf("time object queryx failed: %s\n", err.Error())
		return 0, results, errors.New("no rows found")
	}
	defer rows.Close()
	for rows.Next() {
		var mdl TimeObject
		err := rows.StructScan(&mdl)
		if err != nil {
			logger.Error(err.Error())
			break
		}
		var obj PeriodTimeObject
		err = json.Unmarshal([]byte(mdl.Value), &obj)
		if err != nil {
			logger.Errorf("query abs time list unmarshal failed: %s\n", err.Error())
			break
		}
		obj.Active = obj.IsActive()
		results = append(results, obj)
	}
	return count, results, err
}

func QueryTimeObjectNames(param *goparam.Param) ([]string, error) {
	db := sqldb.GetDB()

	var results []string = []string{}

	queryClause := `select name from object_time`
	queryClause += (param.WhereClause + " order by time_type asc")
	rows, err := db.Queryx(queryClause, param.WhereArgs...)
	if err != nil {
		logger.Errorf("time object queryx failed: %s\n", err.Error())
		return results, errors.New("no rows found")
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		results = append(results, name)
	}
	return results, err
}

func DeleteTimeObject(name string) error {
	m := GetTimeObject(name)
	if m == nil {
		return fmt.Errorf("time object with name %s is not exist", name)
	}
	return m.Remove()
}

func TimeObjectList() ([]TimeObject, error) {
	db := sqldb.GetDB()

	var results []TimeObject = []TimeObject{}

	queryClause := `select name from object_time`
	rows, err := db.Queryx(queryClause)
	if err != nil {
		logger.Errorf("time object queryx failed: %s\n", err.Error())
		return results, errors.New("no rows found")
	}
	defer rows.Close()
	for rows.Next() {
		var mdl TimeObject
		err := rows.StructScan(&mdl)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		results = append(results, mdl)
	}
	return results, err
}
