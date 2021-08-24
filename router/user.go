package router

import (
	userc "github.com/open-cmi/cmmns/controller/user"

	"github.com/gin-gonic/gin"
)

// UserGroup user group router
func UserGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/user")
	{
		g.GET("/", userc.List)
		g.GET("/:id", userc.Get)
		g.GET("/:id/self", userc.GetSelf)
	}
}
