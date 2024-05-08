package licmng

import (
	"errors"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

func ListLicense(query *goparam.Option) (int, []Model, error) {
	db := sqldb.GetConfDB()

	var lics []Model = []Model{}
	countClause := "select count(*) from license"

	whereClause, args := goparam.BuildWhereClause(query)

	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("license list count failed, %s\n", err.Error())
		return 0, lics, errors.New("list count failed")
	}

	queryClause := `select * from license`
	finalClause := goparam.BuildFinalClause(query)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, lics, nil
	}

	for rows.Next() {
		var item Model
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("license struct scan failed %s\n", err.Error())
			break
		}

		lics = append(lics, item)
	}
	return count, lics, err
}

func CreateLicense(req *CreateLicenseRequest) error {
	m := New()
	m.Customer = req.Customer
	m.Prod = req.Prod
	m.Version = req.Version
	m.Modules = req.Modules
	m.ExpireTime = req.ExpireTime
	m.Machine = req.Machine
	return m.Save()
}

func DeleteLicense(id string) error {
	m := Get(id)
	if m == nil {
		return errors.New("license not exist")
	}
	return m.Remove()
}
