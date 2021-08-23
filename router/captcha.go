package router

import (
	captchac "github.com/open-cmi/cmmns/controller/captcha"

	"github.com/gin-gonic/gin"
)

// CaptchaGroup define captcha group
func CaptchaGroup(e *gin.Engine) {
	g := e.Group("/api/common/captcha")
	{
		g.GET("/", captchac.GetID)
		g.GET("/:id", captchac.GetPic)
	}
}
