package time

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

const (
	TimeTypeAbsolute int = iota
	TimeTypePeriod
)

type TimeObject struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	active      bool
	TimeType    int    `json:"time_type" db:"time_type"`
	RefCnt      int    `json:"refcnt" db:"refcnt"`
	Value       string `json:"value" db:"value"`
	isNew       bool   `json:"-"`
}

func (to *TimeObject) Save() error {
	db := sqldb.GetDB()

	if to.isNew {
		// 存储到数据库
		columns := goparam.GetColumn(*to, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into object_time(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s", insertClause)

		_, err := db.NamedExec(insertClause, to)
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
	} else {
		columns := goparam.GetColumn(*to, []string{"name"})

		updates := goparam.GetColumnUpdateNamed(columns)
		updateClause := fmt.Sprintf("update object_time set %s where name=:name", strings.Join(updates, ","))
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.NamedExec(updateClause, to)
		if err != nil {
			logger.Errorf("update object time model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}

	return nil
}

func (to *TimeObject) Remove() error {
	db := sqldb.GetDB()

	deleteClause := "delete from object_time where name=$1"
	_, err := db.Exec(deleteClause, to.Name)
	if err != nil {
		return errors.New("delete object_time failed")
	}
	return nil
}

func (t *TimeObject) IsActive() bool {
	switch t.TimeType {
	case TimeTypeAbsolute:
		var atm AbsoluteTimeObject
		json.Unmarshal([]byte(t.Value), &atm)
		return atm.IsActive()
	case TimeTypePeriod:
		var ptm PeriodTimeObject
		json.Unmarshal([]byte(t.Value), &ptm)
		return ptm.IsActive()
	}

	return false
}

func (m *TimeObject) Ref() error {
	m.RefCnt += 1
	updateClause := "update object_time set refcnt=$1 where name=$2"
	db := sqldb.GetDB()
	_, err := db.Exec(updateClause, m.RefCnt, m.Name)
	if err != nil {
		return errors.New("update object_time refcnt failed")
	}
	return err
}

func (m *TimeObject) Deref() error {
	m.RefCnt -= 1
	updateClause := "update object_time set refcnt=$1 where name=$2"
	db := sqldb.GetDB()
	_, err := db.Exec(updateClause, m.RefCnt, m.Name)
	if err != nil {
		return errors.New("update object_time refcnt failed")
	}
	return err
}

func CreateNewTimeObject(name string, description string, tt int, value string) *TimeObject {
	return &TimeObject{
		Name:        name,
		Description: description,
		TimeType:    tt,
		Value:       value,
		isNew:       true,
	}
}

func GetTimeObject(name string) *TimeObject {
	queryClause := "select * from object_time where name=$1"
	logger.Debugf(queryClause + "\n")
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, name)

	var mdl TimeObject
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return &mdl
}
