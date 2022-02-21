package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jameskeane/bcrypt"
	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/pubsub"
	"github.com/open-cmi/cmmns/essential/storage/sqldb"
	"github.com/open-cmi/cmmns/module/user/msg"
)

type Model struct {
	UserName    string `json:"username" db:"username"`
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"-" db:"password"`
	Role        int    `json:"role" db:"role"`
	Description string `json:"description,omitempty" db:"description"`
	Status      int    `json:"status" db:"status"`
}

// List list func
func List(query *api.Option) (int, []Model, error) {
	db := sqldb.GetDB()

	var users []Model = []Model{}
	countClause := fmt.Sprintf("select count(*) from users")

	whereClause, args := api.BuildWhereClause(query)

	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, users, errors.New("get count failed")
	}

	queryClause := fmt.Sprintf(`select id,username,email,role,description from users`)
	finalClause := api.BuildFinalClause(query)
	queryClause += (whereClause + finalClause)
	rows, err := db.Query(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, users, nil
	}

	for rows.Next() {
		var item Model
		err := rows.Scan(&item.ID, &item.UserName, &item.Email, &item.Role, &item.Description)
		if err != nil {
			break
		}

		users = append(users, item)
	}
	return count, users, err
}

// Get get id
func Get(option *api.Option, field string, value string) (user *Model) {
	columns := api.GetColumn(Model{}, []string{})

	queryClause := fmt.Sprintf(`select %s from users where %s=$1`, strings.Join(columns, ","), field)
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, value)

	var mdl Model
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return &mdl
}

func VerifyPasswordByID(userid string, password string) bool {
	queryclause := fmt.Sprintf("select password from users where id=$1")

	var pass string
	db := sqldb.GetDB()
	row := db.QueryRow(queryclause, userid)
	err := row.Scan(&pass)
	if err != nil {
		// 用户名不存在
		return false
	}
	if !bcrypt.Match(password, pass) {
		// 用户名密码错误
		return false
	}
	return true
}

func ChangePassword(userid string, password string) error {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(password, salt)
	updateClause := fmt.Sprintf("update users set password='%s'", hash)
	db := sqldb.GetDB()
	_, err := db.Exec(updateClause)
	return err
}

// GetByName get by name
func GetByName(name string) (user Model, err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select id,username,email,role,description from users where username='%s'", name)

	db := sqldb.GetDB()
	row := db.QueryRow(queryclause)
	err = row.Scan(&user.ID, &user.UserName, &user.Email, &user.Role, &user.Description)
	if err != nil {
		// 用户名不存在
		return user, errors.New("user not exist")
	}

	return user, nil
}

// Login  user login
func Login(m *msg.LoginMsg) (authuser *Model, err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select id,username,email,password,status from users where username=$1")

	var user Model
	var password string
	db := sqldb.GetDB()
	row := db.QueryRow(queryclause, m.UserName)
	err = row.Scan(&user.ID, &user.UserName, &user.Email, &password, &user.Status)
	if err != nil {
		// 用户名不存在
		return nil, errors.New("username and password not match")
	}

	// 验证密码是否正确， 后续添加salt
	if !bcrypt.Match(m.Password, password) {
		// 用户名密码错误
		return nil, errors.New("username and password not match")
	}
	if user.Status == 0 {
		return nil, errors.New("user has not been activated")
	}
	return &user, nil
}

// Activate activate user
func Activate(username string) error {
	updateClause := fmt.Sprintf("update users set status=1 where username=$1", username)
	db := sqldb.GetDB()
	_, err := db.Exec(updateClause, username)
	if err != nil {
		return errors.New("activate user failed")
	}
	return err
}

// Delete delete user
func DeleteByName(username string) error {
	deleteClause := fmt.Sprintf("delete from users where username=$1")
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, username)
	if err != nil {
		return errors.New("del user failed")
	}
	return err
}

func DeleteByID(id string) error {
	deleteClause := fmt.Sprintf("delete from users where id=$1")
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, id)
	if err != nil {
		return errors.New("del user failed")
	}
	return err
}

// Register register user
func Register(m *msg.RegisterMsg) (err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select username from users where username=$1")

	var un string
	db := sqldb.GetDB()
	row := db.QueryRow(queryclause, m.UserName)
	err = row.Scan(&un)
	if err == nil {
		// 用户名已经被占用
		return errors.New("username has been used")
	}

	queryclause = fmt.Sprintf("select email from users where email=$1")

	var email string
	row = db.QueryRow(queryclause, m.Email)
	err = row.Scan(&email)
	if err == nil {
		// 邮箱已经被占用
		return errors.New("email has been used")
	}

	id := uuid.New()
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(m.Password, salt)
	insertClause := fmt.Sprintf("insert into users(id, username, password, email, description) values($1, $2, $3, $4, $5)")

	_, err = db.Exec(insertClause, id.String(), m.UserName, hash, m.Email, m.Description)
	if err != nil {
		return errors.New("create user failed")
	}
	return nil
}

func Create(m *msg.CreateMsg) (err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select username from users where username=$1 or email=$2")

	var un string
	db := sqldb.GetDB()
	row := db.QueryRow(queryclause, m.UserName, m.Email)
	err = row.Scan(&un)
	if err == nil {
		// 用户名已经被占用
		return errors.New("username or email has been used")
	}

	id := uuid.New()
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(m.Password, salt)
	insertClause := fmt.Sprintf("insert into users(id, username, password, email, status, description) values($1, $2, $3, $4, $5, $6)")

	_, err = db.Exec(insertClause, id.String(), m.UserName, hash, m.Email, 1, m.Description)
	if err != nil {
		return errors.New("create user failed")
	}

	pubsub.Publish(def.EventUserCreate, m.UserName)
	return nil
}
