package router

import (
	"github.com/open-cmi/cmmns/module/secretkey/controller"

	"github.com/gin-gonic/gin"
)

// AuthGroup secretkey auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/secret-key")
	{
		g.GET("/", controller.List)
		g.POST("/", controller.Create)
		g.POST("/multi-delete", controller.MultiDelete)
		g.GET("/:id", controller.Get)
		g.DELETE("/:id", controller.Delete)
		g.PUT("/:id", controller.Edit)
	}
}
