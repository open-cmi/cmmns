package secretkey

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

func GetByName(name string) *Model {
	queryclause := "select * from secret_key where name=$1"

	var model Model
	db := sqldb.GetConfDB()
	row := db.QueryRowx(queryclause, name)
	err := row.StructScan(&model)
	if err != nil {
		// 用户名已经被占用
		logger.Errorf("secret key get %s failed\n", name)
		return nil
	}
	return &model
}

func Get(id string) *Model {
	// 先检查用户名是否存在
	queryclause := "select * from secret_key where id=$1"

	var model Model
	db := sqldb.GetConfDB()
	row := db.QueryRowx(queryclause, id)
	err := row.StructScan(&model)
	if err == nil {
		// 用户名已经被占用
		return &model
	}
	return nil
}

// NameList name list
func NameList() (int, []string, error) {
	db := sqldb.GetConfDB()

	var results []string = []string{}

	countClause := "select count(*) from secret_key"
	row := db.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	queryClause := `select name from secret_key`
	rows, err := db.Query(queryClause)
	if err != nil {
		// 没有的话，也不需要报错
		return count, results, nil
	}

	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

// List list
func List(mo *goparam.Param) (int, []Model, error) {
	db := sqldb.GetConfDB()

	var results []Model = []Model{}

	countClause := "select count(*) from secret_key"
	whereClause, args := goparam.BuildWhereClause(mo)
	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, results, errors.New("get count failed")
	}

	queryClause := `select * from secret_key`
	queryClause += whereClause
	queryClause += goparam.BuildFinalClause(mo)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, results, nil
	}

	for rows.Next() {
		var item Model
		err := rows.StructScan(&item)
		if err != nil {
			break
		}

		results = append(results, item)
	}
	return count, results, err
}

// List list
func MultiDelete(mo *goparam.Param, ids []string) error {
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

	deleteClause := fmt.Sprintf("delete from secret_key where id in %s", list)
	_, err := db.Exec(deleteClause, args...)
	if err != nil {
		return errors.New("delete item failed")
	}
	return nil
}

func Create(mo *goparam.Param, reqMsg *CreateMsg) (m *Model, err error) {
	// 先检查用户名是否存在
	model := Get(reqMsg.Name)
	if model != nil {
		// 用户名已经被占用
		return nil, errors.New("name has been used")
	}

	m = New(reqMsg)
	err = m.Save()

	return m, err
}

func Edit(mo *goparam.Param, id string, reqMsg *EditMsg) error {
	m := Get(id)
	if m == nil {
		return errors.New("item not exist")
	}

	m.Name = reqMsg.Name
	err := m.Save()
	return err
}

func Delete(mo *goparam.Param, id string) error {
	m := Get(id)
	if m == nil {
		return errors.New("item not exist")
	}
	return m.Remove()
}

func CreateByFile(name string, privateFile string) error {
	m := GetByName(name)
	if m != nil {
		return errors.New("name has been used")
	}
	privateKey, err := os.ReadFile(privateFile)
	if err != nil {
		logger.Errorf("read private file failed: %s\n", err.Error())
		return errors.New("read private key file failed")
	}
	publicKey, err := GeneratePublickKey(privateFile)
	if err != nil {
		logger.Errorf("generate public key failed: %s\n", err.Error())
		return errors.New("generate public key failed")
	}
	m = &Model{
		ID:    uuid.NewString(),
		isNew: true,
	}
	m.Name = name
	m.PrivateKey = string(privateKey)
	m.PublicKey = publicKey
	err = m.Save()
	return err
}
