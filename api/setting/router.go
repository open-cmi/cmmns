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
	webserver.RegisterAuthRouter("setting", "/api/common/v3/setting")
	webserver.RegisterAuthAPI("setting", "GET", "/", List)
	webserver.RegisterAuthAPI("setting", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("setting", "PUT", "/:id", Edit)
}
