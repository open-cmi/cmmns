package router

import (
	"github.com/open-cmi/cmmns/module/user/controller"

	"github.com/gin-gonic/gin"
)

// UnauthGroup user nauth group
func UnauthGroup(e *gin.Engine) {
	g := e.Group("/api/common/v3/user")
	{
		g.POST("/login", controller.Login)
		g.POST("/register", controller.Register)
		g.GET("/activate/:code", controller.Activate)
	}
}

// AuthGroup user auth group router
func AuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/user")
	{
		g.GET("/checkauth", controller.CheckAuth)
		g.GET("/", controller.List)
		g.POST("/", controller.Create)
		g.POST("/changepass", controller.ChangePassword)
		g.POST("/logout", controller.Logout)
		g.GET("/:id", controller.Get)
		g.DELETE("/:id", controller.Delete)
	}
}
