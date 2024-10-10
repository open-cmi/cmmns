package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/setting/time"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetNtpSetting(c *gin.Context) {
	m := time.Get()
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": m})
}

func GetTimeZoneList(c *gin.Context) {
	arr, err := time.GetTimeZoneList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": arr})
}

func SetTimeSetting(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req time.SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set ntp setting"), false)
		return
	}
	err := time.SetTimeSetting(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set ntp setting"), false)
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("set ntp setting"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/time-setting/", GetNtpSetting)
	webserver.RegisterAuthAPI("system", "POST", "/time-setting/", SetTimeSetting)

	webserver.RegisterAuthAPI("system", "GET", "/time-setting/tz/", GetTimeZoneList)
}
