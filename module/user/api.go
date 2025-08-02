package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jameskeane/bcrypt"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/pkg/goparam"
)

const UserLoginMaxTried = 5

type QueryFilter struct {
	Username string
}

// List list func
func QueryList(query *goparam.Param, filter *QueryFilter) (int, []UserModel, error) {
	db := sqldb.GetDB()

	startIndex := query.PageParam.Page * query.PageParam.PageSize
	var users []UserModel = []UserModel{}
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
	finalClause := goparam.BuildFinalClause(query, []string{"created_time", "updated_time"})
	queryClause += (whereClause + finalClause)
	rows, err := db.Queryx(queryClause, whereArgs...)
	if err != nil {
		// 没有的话，也不需要报错
		return count, users, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item UserModel
		err := rows.StructScan(&item)
		if err != nil {
			logger.Errorf("user struct scan failed %s\n", err.Error())
			break
		}
		startIndex += 1
		item.Index = startIndex

		users = append(users, item)
	}
	return count, users, err
}

// Login  user login
func Login(m *LoginMsg) (*UserModel, error) {
	// 先检查用户名是否存在
	queryclause := `select * from users where username=$1`
	var user UserModel
	db := sqldb.GetDB()
	row := db.QueryRowx(queryclause, m.UserName)
	err := row.StructScan(&user)
	if err != nil {
		logger.Errorf("Login: struct scan failed: %s\n", err.Error())
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
	// 登陆时提示用户修改密码
	if user.PasswordChangeTime == 0 {
		user.NeedChangePassword = true
	}
	return &user, nil
}

func Create(req *CreateMsg) (err error) {
	if req.Password != req.ConfirmPass {
		return errors.New(i18n.Sprintf("confirm password is not same with password"))
	}
	// 先检查用户名是否存在
	queryclause := "select username from users where username=$1"

	var un string
	db := sqldb.GetDB()
	row := db.QueryRow(queryclause, req.UserName)
	err = row.Scan(&un)
	if err == nil {
		// 用户名已经被占用
		return errors.New(i18n.Sprintf("username is existing"))
	}

	rm := rbac.GetByName(req.Role)
	if rm == nil {
		return errors.New(i18n.Sprintf("role is not existing"))
	}

	id := uuid.New()
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(req.Password, salt)

	user := NewUser()
	user.ID = id.String()
	user.UserName = req.UserName
	user.Password = hash
	user.Email = req.Email
	user.Description = req.Description
	user.Activate = true
	user.Role = req.Role
	user.Status = "offline"
	user.PasswordChangeTime = time.Now().Unix()

	err = user.Save()
	return err
}

func Edit(req *EditMsg) error {
	user := Get(req.ID)
	if user == nil {
		return errors.New("username is not existing")
	}

	rm := rbac.GetByName(req.Role)
	if rm == nil {
		return errors.New(i18n.Sprintf("role is not existing"))
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
	if req.Password != req.ConfirmPassword {
		return errors.New(i18n.Sprintf("password confirmation doesn't match the password"))
	}

	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(req.Password, salt)

	t := time.Now().Unix()
	updateClause := `update users set password=$1,password_change_time=$2 where id=$3`
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
