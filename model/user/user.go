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
	Email    string `json:"email"`
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
	queryclause := fmt.Sprintf("select id,username,email,password,status from users where username='%s'", m.UserName)

	var user User
	var password string
	var status int
	sqldb := db.GetDB()
	row := sqldb.QueryRow(queryclause)
	err = row.Scan(&user.ID, &user.UserName, &user.Email, &password, &status)
	if err != nil {
		// 用户名不存在
		return nil, errors.New("username and password not match")
	}

	// 验证密码是否正确， 后续添加salt
	if !bcrypt.Match(m.Password, password) {
		// 用户名密码错误
		return nil, errors.New("username and password not match")
	}
	if status == 0 {
		return nil, errors.New("user has not been activated")
	}
	return &user, nil
}

// Activate activate user
func Activate(username string) error {
	updateClause := fmt.Sprintf("update users set status=1 where username='%s'", username)
	sqldb := db.GetDB()
	_, err := sqldb.Exec(updateClause)
	if err != nil {
		return errors.New("activate user failed")
	}
	return err
}

// Delete delete user
func Delete(username string) error {
	deleteClause := fmt.Sprintf("delete from users where username='%s'", username)
	sqldb := db.GetDB()
	_, err := sqldb.Exec(deleteClause)
	if err != nil {
		return errors.New("del user failed")
	}
	return err
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
