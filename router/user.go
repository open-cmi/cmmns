package router

import (
	userc "github.com/open-cmi/cmmns/controller/user"

	"github.com/gin-gonic/gin"
)

// UserNauthGroup user nauth group
func UserNauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/user")
	{
		g.POST("/login", userc.Login)
		g.POST("/register", userc.Register)
		g.GET("/activate/:code", userc.Activate)
	}
}

// UserAuthGroup user auth group router
func UserAuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/user")
	{
		g.GET("/checkauth", userc.CheckAuth)
		g.GET("/", userc.List)
		g.POST("/", userc.Create)
		g.POST("/changepass", userc.ChangePassword)
		g.POST("/logout", userc.Logout)
		g.GET("/:id", userc.Get)
		g.DELETE("/:id", userc.Delete)
	}
}
