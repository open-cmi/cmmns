package upgrademng

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/system/upgrademng"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/webserver"
)

func SetUpgradeSetting(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	var req upgrademng.SetSettingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set upgrade setting"), false)
		return
	}
	err := upgrademng.SetUpgradeSetting(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set upgrade setting"), false)
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("set upgrade setting"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func GetUpgradeSetting(c *gin.Context) {
	s := upgrademng.GetUpgradeSetting()
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": s,
	})
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/upgrade-mng/setting/", GetUpgradeSetting)
	webserver.RegisterAuthAPI("system", "POST", "/upgrade-mng/setting/", SetUpgradeSetting)
}
