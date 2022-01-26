package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/storage/db"

	msg "github.com/open-cmi/cmmns/msg/agent"
)

const (
	StateInit = iota
	StateDeploySuccess
	StateDeployFailed
	StateOnline
	StateOffline
)

// Model  model
type Model struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	DevID        string `json:"dev_id" db:"dev_id"`
	Group        string `json:"group" db:"group_name"`
	Address      string `json:"address" db:"address"`
	LocalAddress string `json:"local_address" db:"local_address"`
	Port         int    `json:"port" db:"port"`
	IsLocal      bool   `json:"is_local" db:"is_local"`
	ConnType     string `json:"conn_type" db:"conn_type"`
	User         string `json:"user" db:"username"`
	Password     string `json:"-" db:"password"`
	SecretKey    string `json:"secret_key" db:"secret_key"`
	Location     string `json:"location" db:"location"`
	State        int    `json:"state" db:"state"`
	Reason       string `json:"reason" db:"reason"`
	Description  string `json:"description" db:"description"`
	isNew        bool
}

// GetPassword 获取敏感使用
func (m *Model) GetPassword() {
	// 这里获取密码等敏感信息
}

func (m *Model) Save() error {
	sqldb := db.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := []string{"id", "name", "group_name",
			"address", "port", "conn_type", "username", "password",
			"secret_key", "description", "location"}
		values := db.GetColumnNamed(columns)

		insertClause := fmt.Sprintf("insert into agent(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Logger.Debug("start to exec sql clause: %s", insertClause)

		_, err := sqldb.NamedExec(insertClause, m)
		if err != nil {
			logger.Logger.Error("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
	} else {
		columns := []string{"name", "group_name", "address", "port", "conn_type",
			"username", "password", "secret_key", "description", "location", "state", "reason"}
		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update agent set %s where id=:id", strings.Join(updates, ","))
		logger.Logger.Debug("start to exec sql clause: %s", updateClause)
		_, err := sqldb.NamedExec(updateClause, m)
		if err != nil {
			logger.Logger.Error("update agent model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	sqldb := db.GetDB()

	deleteClause := fmt.Sprintf("delete from agent where id=$1")
	_, err := sqldb.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func New(reqMsg *msg.CreateMsg) (m *Model) {
	return &Model{
		ID:          uuid.NewString(),
		Name:        reqMsg.Name,
		Group:       reqMsg.Group,
		Address:     reqMsg.Address,
		Port:        reqMsg.Port,
		ConnType:    reqMsg.ConnType,
		User:        reqMsg.UserName,
		Password:    reqMsg.Password,
		SecretKey:   reqMsg.SecretKey,
		Location:    reqMsg.Location,
		Description: reqMsg.Description,
		isNew:       true,
	}
}
