package router

import (
	"github.com/open-cmi/cmmns/module/system/controller"

	"github.com/gin-gonic/gin"
)

// AuthGroup system auth group
func AuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/system/")
	{
		g.GET("/status/", controller.List)
		g.GET("/status/:id", controller.Get)

		g.GET("/device/", controller.GetDevID)

		g.POST("/reboot", controller.Reboot)
	}
}
