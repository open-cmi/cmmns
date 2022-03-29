package user

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// UnauthGroup user nauth group
func UnauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/user")
	{
		g.POST("/login", Login)
		g.POST("/register", Register)
		g.GET("/activate/:code", Activate)
	}
}

// AuthGroup user auth group router
func AuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/user")
	{
		g.GET("/checkauth", CheckAuth)
		g.GET("/", List)
		g.POST("/", Create)
		g.POST("/changepass", ChangePassword)
		g.POST("/logout", Logout)
		g.GET("/:id", Get)
		g.DELETE("/:id", Delete)
	}
}

func init() {
	webserver.RegisterAuthAPI("user", AuthGroup)
	webserver.RegisterUnauthAPI("user", UnauthGroup)
}
