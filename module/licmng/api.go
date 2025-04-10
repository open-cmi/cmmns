package licmng

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

func ListLicense(query *goparam.Param) (int, []Model, error) {
	db := sqldb.GetDB()

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
	defer rows.Close()
	for rows.Next() {
		var item Model
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("license struct scan failed %s\n", err.Error())
			break
		}
		item.Serial = GenerateSerial(item.MCode)

		lics = append(lics, item)
	}
	return count, lics, err
}

func CreateLicense(req *CreateLicenseRequest) (*Model, error) {
	if req.Version != "trial" && req.Version != "pro" && req.Version != "enterprise" {
		return nil, fmt.Errorf("not supported version, version should be trial, pro or enterprise")
	}

	if req.Version == "pro" && req.MCode == "" {
		return nil, fmt.Errorf("machine code should not be empty")
	}

	m := New()
	m.Customer = req.Customer
	m.Prod = req.Prod
	m.Version = req.Version
	m.Modules = req.Modules
	m.ExpireTime = req.ExpireTime
	m.MCode = req.MCode
	err := m.Save()
	return m, err
}

func DeleteLicense(id string) error {
	m := Get(id)
	if m == nil {
		return errors.New("license not exist")
	}
	return m.Remove()
}

func GenerateSerial(mcode string) string {
	str := fmt.Sprintf("swapi-%s-mcode", mcode)
	bs64 := base64.StdEncoding.EncodeToString([]byte(str))

	result := md5.Sum([]byte(bs64))

	serial := fmt.Sprintf("0000%02x-%02x%02x%02x%02x%02x-%02x%02x%02x%02x%02x-%02x%02x%02x%02x%02x", result[0], result[1], result[2], result[3], result[4],
		result[5], result[6], result[7], result[8], result[9], result[10], result[11], result[12], result[13], result[14], result[15])
	return serial
}
