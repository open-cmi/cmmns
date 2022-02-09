package router

import (
	"github.com/open-cmi/cmmns/controller/agentgroup"

	"github.com/gin-gonic/gin"
)

// 如果针对某个id的操作，使用/:id/xxx的形式，尽量避免冲突
// 如果针对多个id的操作，可以使用/multi-xxx的形式，类似multi-delete,然后在body中携带参数
// 如果GET有多个/xxx的，可以在group当中就区分开来，将/:id形式的单独放一个组，其他有具体/xxx的放一个组

// AgentGroupNauthGroup agent group nauth group
func AgentGroupNauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent-group-nauth")
	{
		g.GET("/self", agentgroup.List)
		g.GET("/keepalive", agentgroup.List)
	}
}

// AgentGroupAuthGroup agent group auth group router
func AgentGroupAuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent-group")
	{
		g.GET("/", agentgroup.List)
		g.POST("/", agentgroup.Create)
		g.POST("/multi-delete", agentgroup.MultiDelete)
		g.GET("/:id", agentgroup.Get)
		g.DELETE("/:id", agentgroup.Delete)
		g.PUT("/:id", agentgroup.Edit)
	}
}
