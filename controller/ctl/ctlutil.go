package ctl

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func GetUser(c *gin.Context) map[string]interface{} {
	sess, _ := c.Get("session")
	session := sess.(*sessions.Session)
	user, ok := session.Values["user"].(map[string]interface{})
	if ok {
		return user
	}
	return nil
}
