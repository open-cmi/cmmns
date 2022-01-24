package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/controller"
	model "github.com/open-cmi/cmmns/model/auditlog"
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
	user := controller.GetUser(c)
	username, _ := user["username"].(string)

	err := model.InsertLog(ip, username, logtype, action)
	return err
}
