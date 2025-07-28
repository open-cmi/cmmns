package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jameskeane/bcrypt"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

type User struct {
	UserName           string `json:"username" db:"username"`
	ID                 string `json:"id" db:"id"`
	Email              string `json:"email" db:"email"`
	Password           string `json:"-" db:"password"`
	Role               string `json:"role" db:"role"`
	Description        string `json:"description,omitempty" db:"description"`
	Activate           bool   `json:"activate" db:"activate"`
	Status             string `json:"status" db:"status"`
	CreatedTime        int64  `json:"created_time" db:"itime"`
	UpdatedTime        int64  `json:"updated_time" db:"utime"`
	PasswordChangeTime int64  `json:"-" db:"password_change_time"`
	NeedChangePassword bool   `json:"need_change_password"`
	isNew              bool
}

func (u *User) Save() error {
	db := sqldb.GetDB()

	if u.isNew {
		columns := goparam.GetColumn(User{}, []string{})
		values := goparam.GetColumnInsertNamed(columns)

		insertClause := fmt.Sprintf("insert into users(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		_, err := db.NamedExec(insertClause, u)
		if err != nil {
			return errors.New("create user failed")
		}
	} else {
		columns := goparam.GetColumn(User{}, []string{})
		values := goparam.GetColumnUpdateNamed(columns)

		updateClause := fmt.Sprintf("update users set %s where id=:id",
			strings.Join(values, ","))

		_, err := db.NamedExec(updateClause, u)
		if err != nil {
			logger.Errorf("update user failed: %s\n", err.Error())
			return errors.New("update user failed")
		}
	}
	return nil
}

func (m *User) Remove() error {
	deleteClause := "delete from users where id=$1"
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("del user failed")
	}
	return err
}

// Get get id
func Get(id string) (user *User) {
	queryClause := `select * from users where id=$1`
	db := sqldb.GetDB()
	row := db.QueryRowx(queryClause, id)

	var mdl User
	err := row.StructScan(&mdl)
	if err != nil {
		logger.Errorf("user %s not found: %s\n", id, err.Error())
		return nil
	}

	return &mdl
}

func VerifyPasswordByID(userid string, password string) bool {
	queryclause := "select password from users where id=$1"

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

// Activate activate user
func Activate(username string) error {
	updateClause := "update users set status=1 where username=$1"
	db := sqldb.GetDB()
	_, err := db.Exec(updateClause, username)
	if err != nil {
		return errors.New("activate user failed")
	}
	return err
}

// Delete delete user
func DeleteByName(username string) error {
	deleteClause := "delete from users where username=$1"
	db := sqldb.GetDB()
	_, err := db.Exec(deleteClause, username)
	if err != nil {
		return errors.New("del user failed")
	}
	return err
}

func NewUser() *User {
	n := time.Now().Unix()
	return &User{
		isNew:       true,
		CreatedTime: n,
		UpdatedTime: n,
	}
}
