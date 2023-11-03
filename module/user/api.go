package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jameskeane/bcrypt"
	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/pubsub"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

// List list func
func List(query *goparam.Option) (int, []User, error) {
	db := sqldb.GetConfDB()

	var users []User = []User{}
	countClause := "select count(*) from users"

	whereClause, args := goparam.BuildWhereClause(query)

	countClause += whereClause
	row := db.QueryRow(countClause, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("user list count failed, %s\n", err.Error())
		return 0, users, errors.New("list count failed")
	}

	queryClause := `select * from users`
	finalClause := goparam.BuildFinalClause(query)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, args...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, users, nil
	}

	for rows.Next() {
		var item User
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("user struct scan failed %s\n", err.Error())
			break
		}

		users = append(users, item)
	}
	return count, users, err
}

// Login  user login
func Login(m *LoginMsg) (authuser *User, err error) {
	// 先检查用户名是否存在
	queryclause := `select * from users where username=$1`

	var user User
	db := sqldb.GetConfDB()
	row := db.QueryRowx(queryclause, m.UserName)
	err = row.StructScan(&user)
	if err != nil {
		// 用户名不存在
		logger.Warnf("login failed: %s\n", err.Error())
		return nil, errors.New("username and password not match")
	}

	// 验证密码是否正确， 后续添加salt
	if !bcrypt.Match(m.Password, user.Password) {
		// 用户名密码错误
		logger.Warnf("login failed: password is incorrect\n")
		return nil, errors.New("username and password not match")
	}

	if !user.Activate {
		return nil, errors.New("user has not been activated")
	}
	return &user, nil
}

func Create(m *CreateMsg) (err error) {
	// 先检查用户名是否存在
	queryclause := "select username from users where username=$1 or email=$2"

	var un string
	db := sqldb.GetConfDB()
	row := db.QueryRow(queryclause, m.UserName, m.Email)
	err = row.Scan(&un)
	if err == nil {
		// 用户名已经被占用
		return errors.New("username or email is exist")
	}

	id := uuid.New()
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(m.Password, salt)

	user := NewUser()
	user.ID = id.String()
	user.UserName = m.UserName
	user.Password = hash
	user.Email = m.Email
	user.Description = m.Description
	user.Activate = true
	user.Role = m.Role
	user.Status = "offline"

	err = user.Save()
	if err == nil {
		pubsub.Publish(def.EventUserCreate, m.UserName)
	}
	return err
}

// Register register user
func Register(m *RegisterMsg) (err error) {
	// 先检查用户名是否存在
	queryclause := "select username from users where username=$1"

	var un string
	db := sqldb.GetConfDB()
	row := db.QueryRow(queryclause, m.UserName)
	err = row.Scan(&un)
	if err == nil {
		// 用户名已经被占用
		return errors.New("username has been registered")
	}

	queryclause = "select email from users where email=$1"

	var email string
	row = db.QueryRow(queryclause, m.Email)
	err = row.Scan(&email)
	if err == nil {
		// 邮箱已经被占用
		return errors.New("email has been registered")
	}

	id := uuid.New()
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(m.Password, salt)

	user := NewUser()
	user.ID = id.String()
	user.UserName = m.UserName
	user.Password = hash
	user.Email = m.Email
	user.Description = m.Description
	user.Activate = false
	user.Role = "subscriber"
	user.Status = "offline"

	err = user.Save()

	return err
}
