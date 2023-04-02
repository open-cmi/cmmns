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

		g.GET("/network", GetNetwork)

		g.PUT("/network", SetNetwork)
		g.GET("/locale", GetLang)
	}
}

func init() {
	webserver.RegisterAuthRouter("system", "/api/common/v3/system/")
	webserver.RegisterAuthAPI("system", "GET", "/status/", GetStatus)
	webserver.RegisterAuthAPI("system", "GET", "/device/", GetDevID)
	webserver.RegisterAuthAPI("system", "POST", "/reboot/", Reboot)
	webserver.RegisterAuthAPI("system", "POST", "/shutdown/", ShutDown)
	webserver.RegisterAuthAPI("system", "PUT", "/locale", ChangeLang)
	webserver.RegisterAuthAPI("system", "GET", "/network", GetNetwork)
	webserver.RegisterAuthAPI("system", "PUT", "/network", SetNetwork)
	webserver.RegisterAuthAPI("system", "GET", "/locale", GetLang)
}
