package router

import (
	systeminfo "github.com/open-cmi/cmmns/controller/systeminfo"

	"github.com/gin-gonic/gin"
)

// SystemInfoGroup device router
func SystemInfoGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/systeminfo")
	{
		g.GET("/", systeminfo.GetSystemInfo)
	}
}
