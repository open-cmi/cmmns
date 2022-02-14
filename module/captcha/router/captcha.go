package router

import (
	"github.com/open-cmi/cmmns/module/captcha/controller"

	"github.com/gin-gonic/gin"
)

// UnauthGroup define captcha group
func UnauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/captcha")
	{
		g.GET("/", controller.GetID)
		g.GET("/:id", controller.GetPic)
	}
}
