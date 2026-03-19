package sysinfo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/pkg/dev"
)

func GetDeviceInfo(c *gin.Context) {
	mcode := dev.GetDeviceID()
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": mcode})
}

func init() {
	rbac.UnauthAPI("system", "GET", "/mcode/", GetDeviceInfo)
}
