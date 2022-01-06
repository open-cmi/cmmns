package router

import (
	"github.com/open-cmi/cmmns/controller/assist"

	"github.com/gin-gonic/gin"
)

// AssistAuthGroup assist auth group router
func AssistAuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/assist")
	{
		g.GET("/", assist.GetAssist)
		g.POST("/", assist.SetAssist)
	}
}
