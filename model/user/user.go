package user

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jameskeane/bcrypt"
	climsg "github.com/open-cmi/cmmns/climsg/user"
	"github.com/open-cmi/cmmns/db"
)

// User user
type User struct {
	UserName string `json:"username"`
	ID       string `json:"id"`
}

// List list func
func List() ([]User, error) {
	return []User{}, nil
}

// Get get id
func Get(id string) (user *User, err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select id,username from users where id='%s'", id)

	var tmpuser User
	sqldb := db.GetDB()
	row := sqldb.QueryRow(queryclause)
	err = row.Scan(&tmpuser.ID, &tmpuser.UserName)
	if err != nil {
		// 用户名不存在
		return nil, errors.New("user not exist")
	}

	return &tmpuser, nil
}

// Login  user login
func Login(m *climsg.LoginMsg) (authuser *User, err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select id,username,password from users where username='%s'", m.UserName)

	var user User
	var password string
	sqldb := db.GetDB()
	row := sqldb.QueryRow(queryclause)
	err = row.Scan(&user.ID, &user.UserName, &password)
	if err != nil {
		// 用户名不存在
		return nil, errors.New("username and password not match")
	}

	// 验证密码是否正确， 后续添加salt
	if !bcrypt.Match(m.Password, password) {
		// 用户名密码错误
		return nil, errors.New("username and password not match")
	}

	return &user, nil
}

// Register register user
func Register(m *climsg.RegisterMsg) (err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select username from users where username=%s", m.UserName)

	var un string
	sqldb := db.GetDB()
	row := sqldb.QueryRow(queryclause)
	err = row.Scan(&un)
	if err == nil {
		// 用户名已经被占用
		return errors.New("username has been used")
	}

	queryclause = fmt.Sprintf("select email from users where email=%s", m.Email)

	var email string
	row = sqldb.QueryRow(queryclause)
	err = row.Scan(&email)
	if err == nil {
		// 邮箱已经被占用
		return errors.New("email has been used")
	}

	id := uuid.New()
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(m.Password, salt)
	insertClause := fmt.Sprintf("insert into users(id, username, password, email, description) values('%s', '%s', '%s', '%s', '%s')",
		id.String(), m.UserName, hash, m.Email, m.Description)

	_, err = sqldb.Exec(insertClause)
	if err != nil {
		return errors.New("create user failed")
	}
	return nil
}
