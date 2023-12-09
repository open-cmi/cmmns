package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/setting/time"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetNtpSetting(c *gin.Context) {
	m := time.GetNtpSetting()
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": m})
}

func SetNtpSetting(c *gin.Context) {
	var req time.SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := time.SetNtpSetting(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func AdjustNtpTime(c *gin.Context) {
	err := time.Adjust()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthAPI("system-setting", "GET", "/ntp-time-setting/", GetNtpSetting)
	webserver.RegisterAuthAPI("system-setting", "POST", "/ntp-time-setting/", SetNtpSetting)
	webserver.RegisterAuthAPI("system-setting", "POST", "/ntp-time-setting/adjust/", AdjustNtpTime)
}
