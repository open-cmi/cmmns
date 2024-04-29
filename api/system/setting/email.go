package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/setting/email"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetEmail(c *gin.Context) {
	m := email.Get()
	if m == nil {
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": *m})
}

func SetEmail(c *gin.Context) {
	var req email.SetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := email.SetNotifyEmail(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func CheckEmail(c *gin.Context) {
	var req email.SetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := email.CheckEmailSetting(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthAPI("system", "POST", "/email-setting/", SetEmail)
	webserver.RegisterAuthAPI("system", "GET", "/email-setting/", GetEmail)
	webserver.RegisterAuthAPI("system", "POST", "/check-email-setting/", CheckEmail)
}
