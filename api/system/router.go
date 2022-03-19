package system

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// AuthGroup system auth group
func AuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/system/")
	{
		g.GET("/status/", GetStatus)

		g.GET("/device/", GetDevID)

		g.POST("/reboot", Reboot)
		g.POST("/shutdown", ShutDown)
		g.PUT("/locale", ChangeLang)
		g.GET("/locale", GetLang)
	}
}

func init() {
	webserver.RegisterAuthAPI("system", AuthGroup)
}
