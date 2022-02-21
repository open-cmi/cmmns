package router

import (
	"github.com/open-cmi/cmmns/module/agent/controller"

	"github.com/gin-gonic/gin"
)

// AgentNauthGroup agent nauth group
func UnauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent")
	{
		g.GET("/get-job", controller.GetJob)
		g.POST("/report-result", controller.ReportResult)
		g.GET("/keep-alive", controller.KeepAlive)
		g.POST("/register", controller.Register)
	}
}

// AgentAuthGroup agent auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent")
	{
		g.GET("/", controller.List)
		g.POST("/", controller.Create)
		g.DELETE("/:id", controller.Delete)
		g.PUT("/:id", controller.Edit)
		g.POST("/deploy/", controller.Deploy)
	}

	g2 := e.Group("/api/common/v3/master/")
	{
		g2.GET("/setting", controller.GetSetting)
		g2.GET("/auto-master-setting", controller.AutoGetMaster)
		g2.POST("/setting", controller.EditSetting)
	}
}
