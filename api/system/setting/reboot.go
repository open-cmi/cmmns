package setting

import (
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/webserver"
)

func Reboot(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	ah.InsertOperationLog(i18n.Sprintf("reboot system"), true)
	exec.Command("/bin/sh", "-c", "reboot -f").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func ShutDown(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	ah.InsertOperationLog(i18n.Sprintf("shutdown system"), true)

	exec.Command("/bin/sh", "-c", "shutdown -h now").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func init() {
	webserver.RegisterAuthAPI("system", "POST", "/reboot/", Reboot)
	webserver.RegisterAuthAPI("system", "POST", "/shutdown/", ShutDown)
}
