package secretkey

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// AuthGroup secretkey auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/secret-key")
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
	webserver.RegisterAuthRouter("secretkey", "/api/common/v3/secret-key")
	webserver.RegisterAuthAPI("secretkey", "GET", "/", List)
	webserver.RegisterAuthAPI("secretkey", "POST", "/", Create)
	webserver.RegisterAuthAPI("secretkey", "POST", "/multi-delete", MultiDelete)
	webserver.RegisterAuthAPI("secretkey", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("secretkey", "DELETE", "/:id", Delete)
	webserver.RegisterAuthAPI("secretkey", "PUT", "/:id", Edit)
}
