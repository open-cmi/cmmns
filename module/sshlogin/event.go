package sshlogin

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
	"github.com/open-cmi/gobase/essential/ticker"
)

// SSHLoginEvent SSH 登录审计事件
type SSHLoginEvent struct {
	ID         string `json:"id" db:"id"`
	Timestamp  int64  `json:"timestamp" db:"timestamp"`
	Username   string `json:"username" db:"username"`
	RHost      string `json:"rhost" db:"rhost"`
	Port       int    `json:"port" db:"port"`
	Status     string `json:"status" db:"status"`         // "success" | "failed"
	AuthType   string `json:"auth_type" db:"auth_type"`   // "password" | "publickey" | "keyboard-interactive" | "pam"
	IsInvalid  int    `json:"is_invalid" db:"is_invalid"` // 1: 非法用户名探测, 0: 正常
	FailReason string `json:"fail_reason" db:"fail_reason"`
}

// LoginHandler 允许业务模块订阅 SSH 登录事件（如后续扩展告警模块）
type LoginHandler func(event *SSHLoginEvent)

var (
	handlersMu sync.RWMutex
	handlers   []LoginHandler
)

// RegisterLoginHandler 注册 SSH 登录事件回调
func RegisterLoginHandler(h LoginHandler) {
	if h == nil {
		return
	}
	handlersMu.Lock()
	handlers = append(handlers, h)
	handlersMu.Unlock()
}

// SaveSSHLoginEvent 保存 SSH 登录事件到本地数据库
func SaveSSHLoginEvent(event *SSHLoginEvent) error {
	db := sqldb.GetDB()
	if db == nil {
		return fmt.Errorf("local database not initialized")
	}

	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().Unix()
	}

	columns := sqldb.ExtractDBColumns(SSHLoginEvent{}, []string{})
	values := sqldb.GetNamedInsertValuesClauses(columns)
	insertSQL := fmt.Sprintf("INSERT INTO ssh_login_event (%s) VALUES (%s)",
		strings.Join(columns, ","), values)
	_, err := db.NamedExec(insertSQL, event)
	if err != nil {
		logger.Errorf("save ssh_login_event failed: %s", err.Error())
		return err
	}
	return nil
}

// HandleSSHLoginEvent 处理实时捕获的 SSH 登录事件
func HandleSSHLoginEvent(event *SSHLoginEvent) {
	if event == nil {
		return
	}

	if err := SaveSSHLoginEvent(event); err != nil {
		logger.Errorf("failed to save ssh login event locally: %s\n", err.Error())
	}

	if event.Status == "success" {
		logger.Infof("ssh login success: user=%s rhost=%s port=%d auth=%s",
			event.Username, event.RHost, event.Port, event.AuthType)
	} else {
		logger.Infof("ssh login failed: user=%s rhost=%s port=%d auth=%s reason=%s invalid=%d",
			event.Username, event.RHost, event.Port, event.AuthType, event.FailReason, event.IsInvalid)
	}

	handlersMu.RLock()
	list := make([]LoginHandler, len(handlers))
	copy(list, handlers)
	handlersMu.RUnlock()
	for _, h := range list {
		h(event)
	}
}

// QuerySSHLoginEvents 分页查询 SSH 登录日志
func QuerySSHLoginEvents(status string, keyword string, limit int, offset int) (int, []*SSHLoginEvent, error) {
	db := sqldb.GetDB()
	if db == nil {
		return 0, nil, fmt.Errorf("database not initialized")
	}

	var conditions []string
	var args []interface{}

	if status != "" && status != "all" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	if keyword != "" {
		conditions = append(conditions, "(username LIKE ? OR rhost LIKE ?)")
		kw := "%" + keyword + "%"
		args = append(args, kw, kw)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM ssh_login_event %s", whereClause)
	var count int
	if err := db.Get(&count, db.Rebind(countSQL), args...); err != nil {
		return 0, nil, err
	}

	querySQL := fmt.Sprintf(`
		SELECT id, timestamp, username, rhost, port, status, auth_type, is_invalid, fail_reason
		FROM ssh_login_event %s
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	queryArgs := append(args, limit, offset)
	var list []*SSHLoginEvent
	if err := db.Select(&list, db.Rebind(querySQL), queryArgs...); err != nil {
		return 0, nil, err
	}

	return count, list, nil
}

func init() {
	// 每天凌晨 2 点清理超过 30 天的本地历史日志，防止数据库膨胀
	ticker.Register("clean-ssh-login-events-ticker", "0 0 2 * * *", func(s string, i interface{}) {
		db := sqldb.GetDB()
		if db == nil {
			return
		}
		expireTime := time.Now().AddDate(0, 0, -30).Unix()
		deleteSQL := "DELETE FROM ssh_login_event WHERE timestamp < ?"
		if res, err := db.Exec(db.Rebind(deleteSQL), expireTime); err == nil {
			if affected, _ := res.RowsAffected(); affected > 0 {
				logger.Infof("cleaned %d expired ssh login events", affected)
			}
		}
	}, nil)
}
