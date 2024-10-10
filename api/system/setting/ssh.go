package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/module/auditlog"
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
	ah := auditlog.NewAuditHandler(c)
	var req ssh.SetSSHServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set ssh service"), false)
		return
	}
	err := ssh.SetSSHServiceSetting(&req)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("set ssh service"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("set ssh service"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthAPI("system", "POST", "/service/ssh-setting/", SetSSHService)
	webserver.RegisterAuthAPI("system", "GET", "/service/ssh-setting/", GetSSHService)
}
