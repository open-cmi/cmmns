package auditlog

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/sqldb"
)

const (
	LoginType = iota
	OperationType
	SystemType
	WarningType
)

// InsertLog insert audit log
func InsertLog(c *gin.Context, logtype int, action string) error {
	// 获取ip地址
	ip := c.ClientIP()

	// 获取用户
	user := api.GetUser(c)
	if user != nil {
		username, _ := user["username"].(string)

		timestamp := time.Now().Unix()
		id := uuid.New().String()

		// 这里应该通过tunnel来传递，不能直接写数据库，后续再调整
		insertClause := fmt.Sprintf(`insert into audit_log(id, type, ip, username, action, timestamp) 
		values('%s', %d, '%s', '%s', '%s', %d)`, id, logtype, ip, username, action, timestamp)

		//
		db := sqldb.GetDB()
		_, err := db.Exec(insertClause)
		return err
	}
	errMsg := "user not exist"
	logger.Error(errMsg)
	return errors.New(errMsg)
}
