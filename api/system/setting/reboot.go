package setting

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/pkg/goparam"
)

func Reboot(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	userparam := goparam.GetUser(c)
	role := userparam["role"].(string)
	if role != "admin" {
		ah.InsertOperationLog(i18n.Sprintf("reboot system"), false)
		c.JSON(http.StatusForbidden, "")
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("reboot system"), true)
	exec.Command("/bin/sh", "-c", "reboot -f").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func ShutDown(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	userparam := goparam.GetUser(c)
	role := userparam["role"].(string)
	if role != "admin" {
		ah.InsertOperationLog(i18n.Sprintf("shutdown system"), false)
		c.JSON(http.StatusForbidden, "")
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("shutdown system"), true)
	exec.Command("/bin/sh", "-c", "shutdown -h now").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func init() {
	rbac.OptionAuthAPI("system", "POST", "/reboot/", Reboot, rbac.GetInitRoles())
	rbac.OptionAuthAPI("system", "POST", "/shutdown/", ShutDown, rbac.GetInitRoles())
}
