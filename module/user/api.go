package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jameskeane/bcrypt"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/pubsub"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

const UserLoginMaxTried = 5

type QueryFilter struct {
	Username string
}

// List list func
func QueryList(query *goparam.Param, filter *QueryFilter) (int, []User, error) {
	db := sqldb.GetDB()

	var users []User = []User{}
	var paramnum int = 1
	var whereClause string
	var whereArgs []interface{}

	if filter.Username != "" {
		if whereClause != "" {
			whereClause += " and "
		} else {
			whereClause += " and "
		}
		whereClause += fmt.Sprintf(`username like %s`, sqldb.LikePlaceHolder(paramnum))
		whereArgs = append(whereArgs, filter.Username)
		paramnum += 1
	}

	countClause := "select count(*) from users"

	countClause += whereClause
	row := db.QueryRow(countClause, whereArgs...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.Errorf("user list count failed, %s\n", err.Error())
		return 0, users, errors.New("list count failed")
	}

	queryClause := `select * from users`
	finalClause := goparam.BuildFinalClause(query)
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, whereArgs...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, users, nil
	}
	defer rows.Close()
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
	db := sqldb.GetDB()
	row := db.QueryRowx(queryclause, m.UserName)
	err = row.StructScan(&user)
	if err != nil {
		logger.Errorf("struct scan failed while login\n")
		return nil, err
	}

	// 验证密码是否正确， 后续添加salt
	if !bcrypt.Match(m.Password, user.Password) {
		// 用户名密码错误
		logger.Errorf("user %s password match failed\n", m.UserName)
		return nil, errors.New("password match failed")
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
	db := sqldb.GetDB()
	row := db.QueryRow(queryclause, m.UserName, m.Email)
	err = row.Scan(&un)
	if err == nil {
		// 用户名已经被占用
		return errors.New(i18n.Sprint("username or email is existing"))
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
	user.PasswordChangeTime = time.Now().Unix()

	err = user.Save()
	if err == nil {
		pubsub.Publish(pubsub.EventUserCreate, m.UserName)
	}
	return err
}

// Register register user
func Register(m *RegisterMsg) (err error) {
	// 先检查用户名是否存在
	queryclause := "select username from users where username=$1"

	var un string
	db := sqldb.GetDB()
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

func Edit(req *EditMsg) error {
	user := Get(req.ID)
	if user == nil {
		return errors.New("user is is not existing")
	}
	user.Email = req.Email
	user.Role = req.Role
	user.Description = req.Description
	user.UserName = req.Username
	err := user.Save()
	return err
}

func ChangePassword(userid string, password string) error {
	u := Get(userid)
	if u == nil {
		return errors.New("users is not existing")
	}

	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(password, salt)
	u.Password = hash
	u.PasswordChangeTime = time.Now().Unix()
	return u.Save()
}

func ResetPasswd(req *ResetPasswdRequest) error {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(req.Password, salt)

	t := time.Now().Unix()
	updateClause := `update users set password=$1,password_change_time=$2 andwhere id=$3`
	db := sqldb.GetDB()
	_, err := db.Exec(updateClause, hash, t, req.ID)
	return err
}

func Delete(id string) error {
	u := Get(id)
	if u == nil {
		return errors.New("user is not existing")
	}
	err := u.Remove()
	return err
}
