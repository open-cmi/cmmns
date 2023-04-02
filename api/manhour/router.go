package manhour

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// AuthGroup manhour auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/manhour")
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
	webserver.RegisterAuthRouter("manhour", "/api/common/v3/manhour")
	webserver.RegisterAuthAPI("manhour", "GET", "/", List)
	webserver.RegisterAuthAPI("manhour", "POST", "/", Create)
	webserver.RegisterAuthAPI("manhour", "POST", "/multi-delete", MultiDelete)
	webserver.RegisterAuthAPI("manhour", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("manhour", "DELETE", "/:id", Delete)
	webserver.RegisterAuthAPI("manhour", "PUT", "/:id", Edit)
}
