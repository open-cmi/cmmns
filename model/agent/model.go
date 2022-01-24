package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/logger"
	"github.com/open-cmi/cmmns/storage/db"
	"github.com/open-cmi/goutils/logutil"

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
	ID          string `json:"id"`
	Name        string `json:"name"`
	DeviceID    string `json:"dev_id"`
	Group       string `json:"group"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	IsLocal     bool   `json:"is_local"`
	ConnType    string `json:"conn_type"`
	User        string `json:"user"`
	Password    string `json:"-"`
	SecretKey   string `json:"secret_key"`
	Location    string `json:"location"`
	State       int    `json:"state"`
	Reason      string `json:"reason"`
	Description string `json:"description"`
	isNew       bool
}

// GetPassword 获取敏感使用
func (m *Model) GetPassword() {
	// 这里获取密码等敏感信息
}

func (m *Model) Save() error {
	dbsql := db.GetDB()

	if m.isNew {
		// 存储到数据库
		columns := []string{"id", "name", "group_name",
			"address", "port", "conn_type", "username", "password",
			"secret_key", "description", "location"}
		var values []string = []string{}
		for index, _ := range columns {
			seq := fmt.Sprintf(`$%d`, index+1)
			values = append(values, seq)
		}

		insertClause := fmt.Sprintf("insert into agent(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))
		logger.Logger.Printf(logger.Info, "%s", insertClause)
		_, err := dbsql.Exec(insertClause, m.ID, m.Name, m.Group, m.Address, m.Port, m.ConnType, m.User,
			m.Password, m.SecretKey, m.Description, m.Location)
		if err != nil {
			return errors.New("create model failed")
		}
	} else {
		columns := []string{"name", "group_name",
			"address", "port", "conn_type", "username", "password",
			"secret_key", "description", "location", "state", "reason"}
		var updates []string = []string{}
		for index, column := range columns {
			updates = append(updates, fmt.Sprintf(`%s=$%d`, column, index+1))
		}
		updateClause := fmt.Sprintf("update agent set %s where id=$%d", strings.Join(updates, ","), len(columns)+1)
		_, err := dbsql.Exec(updateClause, m.Name, m.Group, m.Address, m.Port, m.ConnType, m.User,
			m.Password, m.SecretKey, m.Description, m.Location, m.State, m.Reason, m.ID)
		if err != nil {
			logger.Logger.Printf(logutil.Error, "update agent model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}

	return nil
}

func (m *Model) Remove() error {
	dbsql := db.GetDB()

	deleteClause := fmt.Sprintf("delete from agent where id=$1")
	_, err := dbsql.Exec(deleteClause, m.ID)
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
