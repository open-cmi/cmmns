package router

import (
	"github.com/open-cmi/cmmns/controller/agent"

	"github.com/gin-gonic/gin"
)

// AgentNauthGroup agent nauth group
func AgentNauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent")
	{
		g.GET("/self", agent.GetSelfTask)
		g.GET("/keepalive", agent.KeepAlive)
	}
}

// AgentAuthGroup agent auth group router
func AgentAuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent")
	{
		g.GET("/", agent.List)
		g.POST("/", agent.CreateAgent)
		g.DELETE("/:id", agent.DelAgent)
		g.POST("/deploy/", agent.DeployAgent)
	}
}
