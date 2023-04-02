package template

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// AuthGroup template auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/template")
	{
		g.GET("/", List)
		g.POST("/", Create)
		g.POST("/multi-delete", MultiDelete)
		g.GET("/:id", Get)
		g.DELETE("/:id", Delete)
		g.PUT("/:id", Edit)
	}
}

func init() {
	webserver.RegisterAuthRouter("template", "/api/common/v3/template")
	webserver.RegisterAuthAPI("template", "GET", "/", List)
	webserver.RegisterAuthAPI("template", "POST", "/", Create)
	webserver.RegisterAuthAPI("template", "POST", "/multi-delete", MultiDelete)
	webserver.RegisterAuthAPI("template", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("template", "DELETE", "/:id", Delete)
	webserver.RegisterAuthAPI("template", "PUT", "/:id", Edit)
}
