package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/storage/sqldb"
)

const (
	StateInit          = iota
	StateDeploySuccess // 部署成功
	StateDeployFailed  // 部署失败
	StateOnline        // 在线
	StateOffline       // 离线
	StateDeny          // 被用户禁用
)

// Model  model
type Model struct {
	ID           string `json:"id" db:"id"`
	CreatedTime  int64  `json:"created_time" db:"created_time"`
	UpdatedTime  int64  `json:"updated_time" db:"updated_time"`
	HostName     string `json:"hostname" db:"hostname"`
	DevID        string `json:"dev_id" db:"dev_id"`
	Group        string `json:"group" db:"group_name"`
	Address      string `json:"address" db:"address"`
	LocalAddress string `json:"local_address" db:"local_address"`
	Port         int    `json:"port" db:"port"`
	ConnType     string `json:"conn_type" db:"conn_type"`
	UserName     string `json:"username" db:"username"`
	Passwd       string `json:"-" db:"passwd"`
	SecretKey    string `json:"secret_key" db:"secret_key"`
	State        int    `json:"state" db:"state"`
	Description  string `json:"description" db:"description"`
	isNew        bool
}

// GetPasswd 获取敏感使用
func (m *Model) GetPasswd() {
	// 这里获取密码等敏感信息
}

func (m *Model) Save() error {
	db := sqldb.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := api.GetColumn(*m, []string{})
		values := api.GetColumnNamed(columns)

		insertClause := fmt.Sprintf("insert into agent(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s", insertClause)

		_, err := db.NamedExec(insertClause, m)
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
	} else {
		columns := api.GetColumn(*m, []string{"id", "created_time"})

		m.UpdatedTime = time.Now().Unix()
		var updates []string = []string{}
		for _, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=:%s`, column, column))
		}
		updateClause := fmt.Sprintf("update agent set %s where id=:id", strings.Join(updates, ","))
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.NamedExec(updateClause, m)
		if err != nil {
			logger.Errorf("update agent model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	db := sqldb.GetDB()

	deleteClause := "delete from agent where id=$1"
	_, err := db.Exec(deleteClause, m.ID)
	if err != nil {
		return errors.New("delete model failed")
	}
	return nil
}

func New() (m *Model) {
	now := time.Now().Unix()
	return &Model{
		ID:          uuid.NewString(),
		CreatedTime: now,
		UpdatedTime: now,
		isNew:       true,
	}
}
