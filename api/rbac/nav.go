package rbac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/pkg/goparam"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetNavMenu(c *gin.Context) {
	param := goparam.ParseParams(c)

	menu := rbac.GetNavMenu(param.Role)
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": menu,
	})
}

func init() {
	webserver.RegisterAuthAPI("rbac", "GET", "/nav-menu/", GetNavMenu)
}
