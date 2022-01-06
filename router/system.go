package router

import (
	"github.com/open-cmi/cmmns/controller/system"

	"github.com/gin-gonic/gin"
)

// SystemInfoGroup device router
func SystemGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/system")
	{
		g.GET("/info", system.GetSystemInfo)
		g.POST("/reboot", system.Reboot)
	}
}
