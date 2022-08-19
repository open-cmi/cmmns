package assist

import (
	"github.com/open-cmi/cmmns/service/webserver"

	"github.com/gin-gonic/gin"
)

// AuthGroup assist auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/assist")
	{
		g.GET("/", GetAssist)
		g.POST("/", SetAssist)
	}
}

func init() {
	webserver.RegisterAuthRouter("assist", AuthGroup)
}
