package scheduler

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// AuthGroup scheduler auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/scheduler")
	{
		g.GET("/", List)
		g.POST("/multi-delete", MultiDelete)
		g.GET("/:id", Get)
		g.DELETE("/:id", Delete)
	}
}

func init() {
	webserver.RegisterAuthRouter("scheduler", "/api/common/v3/scheduler")
	webserver.RegisterAuthAPI("scheduler", "GET", "/", List)
	webserver.RegisterAuthAPI("scheduler", "POST", "/multi-delete", MultiDelete)
	webserver.RegisterAuthAPI("scheduler", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("scheduler", "DELETE", "/:id", Delete)
}
