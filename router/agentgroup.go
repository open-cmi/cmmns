package router

import (
	"github.com/open-cmi/cmmns/controller/agentgroup"

	"github.com/gin-gonic/gin"
)

// AgentGroupGroup agent group router
func AgentGroupGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agentgroup")
	{
		g.GET("/", agentgroup.List)
	}
}
