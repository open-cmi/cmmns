package router

import (
	"github.com/open-cmi/cmmns/controller/master"

	"github.com/gin-gonic/gin"
)

// MasterAuthGroup master auth group router
func MasterAuthGroup(e *gin.Engine) {

	g2 := e.Group("/api/common/v3/master/")
	{
		g2.GET("/setting", master.GetSetting)
		g2.GET("/auto-get-master", master.AutoGetMaster)
		g2.POST("/setting", master.EditSetting)
	}
}
