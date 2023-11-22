package sysinfo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/system"
)

func GetBasicSystemInfo(c *gin.Context) {
	sysinfo, err := system.GetBasicSystemInfo()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": sysinfo})
}
