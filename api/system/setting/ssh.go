package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/setting/service/ssh"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetSSHService(c *gin.Context) {
	m := ssh.GetSSHServiceSetting()
	if m == nil {
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": *m})
}

func SetSSHService(c *gin.Context) {
	var req ssh.SetSSHServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := ssh.SetSSHServiceSetting(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthAPI("system", "POST", "/service/ssh-setting/", SetSSHService)
	webserver.RegisterAuthAPI("system", "GET", "/service/ssh-setting/", GetSSHService)
}
