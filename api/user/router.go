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
	webserver.RegisterAuthRouter("user", "/api/common/v3/user")
	webserver.RegisterAuthAPI("user", "GET", "/checkauth", CheckAuth)
	webserver.RegisterAuthAPI("user", "GET", "/", List)
	webserver.RegisterAuthAPI("user", "POST", "/", Create)
	webserver.RegisterAuthAPI("user", "POST", "/changepass", ChangePassword)
	webserver.RegisterAuthAPI("user", "POST", "/logout", Logout)
	webserver.RegisterUnauthAPI("user", "POST", "/jwt-token/", CreateToken)
	webserver.RegisterAuthAPI("user", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("user", "DELETE", "/:id", Delete)

	webserver.RegisterUnauthRouter("user", "/api/common/v3/user")
	webserver.RegisterUnauthAPI("user", "POST", "/login", Login)
	webserver.RegisterUnauthAPI("user", "POST", "/register", Register)
	webserver.RegisterUnauthAPI("user", "GET", "/activate/:code", Activate)
}
