package sysinfo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/webserver"
	"github.com/open-cmi/cmmns/module/system"
)

func GetSystemStatusInfo(c *gin.Context) {
	sysinfo, err := system.GetSystemStatusInfo()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": sysinfo})
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/system-status/", GetSystemStatusInfo)
}
