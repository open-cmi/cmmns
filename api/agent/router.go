package agent

import (
	"github.com/gin-gonic/gin"
)

// AgentNauthGroup agent nauth group
func UnauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent")
	{
		g.GET("/get-job", GetJob)
		g.POST("/report-result", ReportResult)
		g.GET("/keep-alive", KeepAlive)
		g.POST("/register", Register)
	}
}

// AgentAuthGroup agent auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/agent")
	{
		g.GET("/", List)
		g.POST("/", Create)
		g.DELETE("/:id", Delete)
		g.PUT("/:id", Edit)
		g.POST("/deploy/", Deploy)
	}

	g2 := e.Group("/api/common/v3/master/")
	{
		g2.GET("/setting", GetSetting)
		g2.GET("/auto-master-setting", AutoGetMaster)
		g2.POST("/setting", EditSetting)
	}
}
