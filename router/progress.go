package router

import (
	"github.com/open-cmi/cmmns/controller/progress"

	"github.com/gin-gonic/gin"
)

// AgentAuthGroup agent auth group router
func ProgressAuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/progress")
	{
		g.GET("/:id", progress.Get)
	}
}
