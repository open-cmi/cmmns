package captcha

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// UnauthGroup define captcha group
func UnauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/captcha")
	{
		g.GET("/", GetID)
		g.GET("/:id", GetPic)
	}
}

func init() {
	webserver.RegisterUnauthRouter("captcha", "/api/common/v3/captcha")
	webserver.RegisterUnauthAPI("captcha", "GET", "/", GetID)
	webserver.RegisterUnauthAPI("captcha", "GET", "/:id", GetPic)
}
