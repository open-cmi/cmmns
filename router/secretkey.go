package router

import (
	"github.com/open-cmi/cmmns/controller/secretkey"

	"github.com/gin-gonic/gin"
)

// SecretKeyAuthGroup secretkey auth group router
func SecretKeyAuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/secret-key")
	{
		g.GET("/", secretkey.List)
		g.POST("/", secretkey.Create)
		g.POST("/multi-delete", secretkey.MultiDelete)
		g.GET("/:id", secretkey.Get)
		g.DELETE("/:id", secretkey.Delete)
		g.PUT("/:id", secretkey.Edit)
	}
}
