package router

import (
	"github.com/open-cmi/cmmns/module/system/controller"

	"github.com/gin-gonic/gin"
)

// AuthGroup system auth group
func AuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/system")
	{
		g.GET("/", controller.List)
		g.GET("/:id", controller.Get)

		g.POST("/reboot", controller.Reboot)
	}
}
