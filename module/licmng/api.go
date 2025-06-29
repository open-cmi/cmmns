package licmng

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

func ListLicense(query *goparam.Param) (int, []LicenseModel, error) {
	db := sqldb.GetDB()

	var lics []LicenseModel = []LicenseModel{}
	countClause := "select count(*) from license"

	countClause += query.WhereClause
	row := db.QueryRow(countClause, query.WhereArgs...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("license list count failed, %s\n", err.Error())
		return 0, lics, errors.New("list count failed")
	}

	queryClause := `select * from license`
	finalClause := goparam.BuildFinalClause(query)
	queryClause += (query.WhereClause + finalClause)
	rows, err := db.Queryx(queryClause, query.WhereArgs...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, lics, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item LicenseModel
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("license struct scan failed %s\n", err.Error())
			break
		}
		item.Serial = GenerateSerial(item.Version, item.Model, item.MCode, item.ExpireTime)

		lics = append(lics, item)
	}
	return count, lics, err
}

func CreateLicense(req *CreateLicenseRequest, username string) (*LicenseModel, error) {
	if req.Version != "trial" && req.Version != "pro" && req.Version != "enterprise" {
		return nil, fmt.Errorf("not supported version")
	}

	if req.Version == "pro" && req.MCode == "" {
		return nil, fmt.Errorf("machine code should not be empty")
	}

	if req.Version == "trial" && (req.ValidPeriod < 1 || req.ValidPeriod > 6) {
		return nil, fmt.Errorf("invalid period")
	}

	if req.Model != "standard" && req.Model != "mini" {
		return nil, fmt.Errorf("invalid model")
	}

	p := time.Now()
	if req.Version == "trial" {
		p = p.AddDate(0, req.ValidPeriod, 0)
	} else {
		p = p.AddDate(0, 999*12, 0)
	}

	m := New()
	m.Customer = req.Customer
	m.Prod = req.Prod
	m.Version = req.Version
	m.Modules = req.Modules
	m.ExpireTime = p.Unix()
	m.MCode = req.MCode
	m.Model = req.Model
	m.Username = username
	err := m.Save()
	return m, err
}

func DeleteLicense(id string) error {
	m := Get(id)
	if m == nil {
		return errors.New("license is not existing")
	}
	return m.Remove()
}

func GenerateGeneralSerial(mcode string) string {
	str := fmt.Sprintf("swapi-%s-mcode", mcode)
	bs64 := base64.StdEncoding.EncodeToString([]byte(str))

	result := md5.Sum([]byte(bs64))

	serial := fmt.Sprintf("0000%02x-%02x%02x%02x%02x%02x-%02x%02x%02x%02x%02x-%02x%02x%02x%02x%02x", result[0], result[1], result[2], result[3], result[4],
		result[5], result[6], result[7], result[8], result[9], result[10], result[11], result[12], result[13], result[14], result[15])
	return serial
}

func GenerateSerial(version string, model string, mcode string, expire int64) string {
	if version == "enterprise" {
		mcode = ""
	}
	str := fmt.Sprintf("%s-%s-%s", version, model, mcode)
	bs64 := base64.StdEncoding.EncodeToString([]byte(str))

	result := md5.Sum([]byte(bs64))

	serial := fmt.Sprintf("000%02x-%02x%02x%02x%02x%02x-%02x%02x%02x%02x%02x-%02x%02x%02x%02x%02x-%x", result[0], result[1], result[2], result[3], result[4],
		result[5], result[6], result[7], result[8], result[9], result[10], result[11], result[12], result[13], result[14], result[15], expire)
	return serial
}

func IsEnable() bool {
	return gConf.Enable
}
