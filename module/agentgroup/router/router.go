package router

import (
	"github.com/open-cmi/cmmns/module/agentgroup/controller"

	"github.com/gin-gonic/gin"
)

// AuthGroup agent group auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent-group")
	{
		g.GET("/", controller.List)
		g.POST("/", controller.Create)
		g.POST("/multi-delete", controller.MultiDelete)
		g.GET("/:id", controller.Get)
		g.DELETE("/:id", controller.Delete)
		g.PUT("/:id", controller.Edit)
	}
}
