package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/essential/webserver"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/setting/service/web"
)

func GetServicePort(c *gin.Context) {
	m := web.GetServicePort()
	if m == nil {
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": *m})
}

func SetServicePort(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	var req web.SetServicePortRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set web service port"), false)
		return
	}
	err := web.SetServicePort(&req)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("set web service port"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("set web service port"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthAPI("system", "POST", "/service/web-setting/", SetServicePort)
	webserver.RegisterAuthAPI("system", "GET", "/service/web-setting/", GetServicePort)
}
