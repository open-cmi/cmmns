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
		g.POST("/", agent.Create)
		g.DELETE("/:id", agent.Delete)
		g.POST("/deploy/", agent.DeployAgent)

	}

	g2 := e.Group("/api/common/v3/agent-setting")
	{
		g2.GET("/", agent.GetSetting)
		g2.GET("/auto-get-master", agent.AutoGetMaster)
		g2.POST("/", agent.EditSetting)
	}
}
