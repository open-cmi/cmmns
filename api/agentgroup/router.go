package agentgroup

import (
	"github.com/open-cmi/cmmns/service/webserver"

	"github.com/gin-gonic/gin"
)

// AuthGroup agent group auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent-group")
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
	webserver.RegisterAuthRouter("agentgroup", "/api/common/v3/agent-group")
	webserver.RegisterAuthAPI("agentgroup", "GET", "/", List)
	webserver.RegisterAuthAPI("agentgroup", "POST", "/", Create)
	webserver.RegisterAuthAPI("agentgroup", "POST", "/multi-delete", MultiDelete)
	webserver.RegisterAuthAPI("agentgroup", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("agentgroup", "DELETE", "/:id", Delete)
	webserver.RegisterAuthAPI("agentgroup", "PUT", "/:id", Edit)
}
