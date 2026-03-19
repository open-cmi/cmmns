package sysinfo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/system/hostname"
	"github.com/open-cmi/gobase/essential/i18n"
)

func SetHostname(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req hostname.SetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set host name"), false)
		return
	}
	err := hostname.Set(&req)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("set host name"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("set host name"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func init() {
	rbac.OptionAuthAPI("system", "POST", "/hostname/", SetHostname, rbac.GetInitRoles())
}
