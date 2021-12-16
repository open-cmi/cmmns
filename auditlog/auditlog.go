package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	model "github.com/open-cmi/cmmns/model/auditlog"
	"github.com/open-cmi/cmmns/model/user"
)

func InsertWebLog(c *gin.Context, action string) error {
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

	err := model.InsertLog(ip, username, model.LoginType, action)
	return err
}
