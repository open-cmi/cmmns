package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	model "github.com/open-cmi/cmmns/model/auditlog"
	"github.com/open-cmi/cmmns/model/user"
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
		userinfo, ok := session.Values["user"].(*user.BasicInfo)
		if ok {
			username = userinfo.UserName
		}
	}

	err := model.InsertLog(ip, username, logtype, action)
	return err
}
