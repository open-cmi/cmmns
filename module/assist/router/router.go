package router

import (
	"github.com/open-cmi/cmmns/module/assist/controller"

	"github.com/gin-gonic/gin"
)

// AuthGroup assist auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/assist")
	{
		g.GET("/", controller.GetAssist)
		g.POST("/", controller.SetAssist)
	}
}
