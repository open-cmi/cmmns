package ssh

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/pkg/shell"
	"github.com/open-cmi/cmmns/pkg/systemctl"
)

type SSHServiceModel struct {
	Port   int  `json:"port"`
	Enable bool `json:"enable"`
	isNew  bool
}

func (m *SSHServiceModel) Key() string {
	return "ssh-service-setting"
}

func (m *SSHServiceModel) Value() string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (m *SSHServiceModel) Save() error {
	db := sqldb.GetConfDB()

	if m.isNew {
		// 存储到数据库
		columns := []string{"key", "value"}
		values := []string{"$1", "$2"}

		insertClause := fmt.Sprintf("insert into k_v_table(%s) values(%s)",
			strings.Join(columns, ","), strings.Join(values, ","))

		logger.Debugf("start to exec sql clause: %s\n", insertClause)
		_, err := db.Exec(insertClause, m.Key(), m.Value())
		if err != nil {
			logger.Errorf("create model failed: %s", err.Error())
			return errors.New("create model failed")
		}
		m.isNew = false
	} else {
		updateClause := "update k_v_table set value=$1 where key=$2"
		logger.Debugf("start to exec sql clause: %s", updateClause)
		_, err := db.Exec(updateClause, m.Value(), m.Key())
		if err != nil {
			logger.Errorf("update model failed: %s", err.Error())
			return errors.New("update model failed")
		}
	}
	return nil
}

func GetSSHServiceModel() *SSHServiceModel {
	queryClause := "select value from k_v_table where key=$1"
	db := sqldb.GetConfDB()
	s := &SSHServiceModel{}
	row := db.QueryRowx(queryClause, s.Key())
	if row == nil {
		return nil
	}
	var v string
	err := row.Scan(&v)
	if err != nil {
		logger.Errorf("row scan failed: %s\n", err.Error())
		return nil
	}
	var m SSHServiceModel
	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		logger.Errorf("model unmarshal failed: %s\n", err.Error())
		return nil
	}
	return &m
}

func NewSSHServiceModel() *SSHServiceModel {
	return &SSHServiceModel{
		isNew: true,
	}
}

type SetSSHServiceRequest struct {
	Port   int  `json:"port"`
	Enable bool `json:"enable"`
}

func SetSSHServiceSetting(req *SetSSHServiceRequest) error {
	m := GetSSHServiceModel()
	if m == nil {
		m = NewSSHServiceModel()
		m.Enable = true
		m.Port = 22
	}
	oen := m.Enable
	if m.Port != req.Port {
		shell.Execute(`sed -i 's/#Port /Port /g' /etc/ssh/sshd_config`)
		oldStr := fmt.Sprintf("Port %d", m.Port)
		newStr := fmt.Sprintf("Port %d", req.Port)
		cmd := fmt.Sprintf(`sed -i 's/%s/%s/g' /etc/ssh/sshd_config`, oldStr, newStr)
		err := shell.Execute(cmd)
		if err != nil {
			return err
		}
	}
	m.Port = req.Port
	m.Enable = req.Enable
	err := m.Save()
	if err != nil {
		return err
	}
	if oen && req.Enable {
		err = systemctl.RestartService("ssh")
	} else if oen && !req.Enable {
		err = systemctl.StopService("ssh")
	} else if !oen && req.Enable {
		err = systemctl.StartService("ssh")
	}
	return err
}

func GetSSHServiceSetting() *SSHServiceModel {
	m := GetSSHServiceModel()
	if m == nil {
		m = NewSSHServiceModel()
		m.Enable = true
		m.Port = 22
	}
	return m
}
