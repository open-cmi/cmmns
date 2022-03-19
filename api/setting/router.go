package setting

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// AuthGroup setting auth group router
func AuthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/setting")
	{
		g.GET("/", List)
		g.GET("/:id", Get)
		g.PUT("/:id", Edit)
	}
}

func init() {
	webserver.RegisterAuthAPI("setting", AuthGroup)
}
