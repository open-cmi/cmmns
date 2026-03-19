package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/setting/email"
	"github.com/open-cmi/gobase/essential/i18n"
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
	ah := auditlog.NewAuditHandler(c)
	var req email.SetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set email setting"), false)
		return
	}
	err := email.SetNotifyEmail(&req)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("set email setting"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("set email setting"), true)

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
	rbac.OptionAuthAPI("system", "POST", "/email-setting/", SetEmail, rbac.GetInitRoles())
	rbac.OptionAuthAPI("system", "GET", "/email-setting/", GetEmail, rbac.GetInitRoles())
	rbac.OptionAuthAPI("system", "POST", "/check-email-setting/", CheckEmail, rbac.GetInitRoles())
}
