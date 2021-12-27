package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	model "github.com/open-cmi/cmmns/model/auditlog"
)

const (
	LoginType     = 0
	OperationType = 1
	SystemType    = 2
	WarningType   = 3
)

func InsertLog(c *gin.Context, logtype int, action string) error {
	// 获取ip地址
	ip := c.ClientIP()

	// 获取用户
	var username string = "unknown"
	s, ok := c.Get("session")
	session, ok := s.(*sessions.Session)

	if ok {
		userinfo, ok := session.Values["user"].(map[string]interface{})
		if ok {
			username = userinfo["username"].(string)
		}
	}

	err := model.InsertLog(ip, username, logtype, action)
	return err
}
