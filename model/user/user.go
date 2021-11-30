package user

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jameskeane/bcrypt"
	"github.com/open-cmi/cmmns/db"
	msg "github.com/open-cmi/cmmns/msg/user"
)

// User user
type User struct {
	UserName string `json:"username"`
	ID       string `json:"id"`
	Email    string `json:"email"`
}

// Model user model
type Model struct {
	UserName    string `json:"username"`
	ID          string `json:"id"`
	Email       string `json:"email"`
	Role        int    `json:"role"`
	Description string `json:"description"`
}

// List list func
func List() (int, []Model, error) {
	dbsql := db.GetDB()

	var users []Model = []Model{}

	countClause := fmt.Sprintf("select count(*) from users")
	row := dbsql.QueryRow(countClause)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, users, errors.New("get count failed")
	}

	queryClause := fmt.Sprintf(`select id,username,email,role,description from users`)

	rows, err := dbsql.Query(queryClause)
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
func Get(id string) (user *User, err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select id,username,email from users where id='%s'", id)

	var tmpuser User
	sqldb := db.GetDB()
	row := sqldb.QueryRow(queryclause)
	err = row.Scan(&tmpuser.ID, &tmpuser.UserName, &tmpuser.Email)
	if err != nil {
		// 用户名不存在
		return nil, errors.New("user not exist")
	}

	return &tmpuser, nil
}

// GetByName get by name
func GetByName(name string) (user User, err error) {
	// 先检查用户名是否存在
	queryclause := fmt.Sprintf("select id,username,email from users where username='%s'", name)

	sqldb := db.GetDB()
	row := sqldb.QueryRow(queryclause)
	err = row.Scan(&user.ID, &user.UserName, &user.Email)
	if err != nil {
		// 用户名不存在
		return user, errors.New("user not exist")
	}

	return user, nil
}

// Login  user login
func Login(m *msg.LoginMsg) (authuser *User, err error) {
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
func Register(m *msg.RegisterMsg) (err error) {
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
