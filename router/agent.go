package router

import (
	"github.com/open-cmi/cmmns/controller/agent"

	"github.com/gin-gonic/gin"
)

// AgentNauthGroup agent nauth group
func AgentNauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent")
	{
		g.GET("/self", agent.GetJob)
		g.GET("/keepalive", agent.KeepAlive)
		g.POST("/register", agent.Register)
	}
}

// AgentAuthGroup agent auth group router
func AgentAuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent")
	{
		g.GET("/", agent.List)
		g.POST("/", agent.Create)
		g.DELETE("/:id", agent.Delete)
		g.PUT("/:id", agent.Edit)
		g.POST("/deploy/", agent.Deploy)
	}
}
