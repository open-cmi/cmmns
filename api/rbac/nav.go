package rbac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/gobase/essential/webserver"
	"github.com/open-cmi/gobase/pkg/goparam"
)

func GetRoleMenus(c *gin.Context) {
	param := goparam.ParseParams(c)

	menu := rbac.GetRoleMenus(param.Role)
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": menu,
	})
}

func init() {
	webserver.RegisterMustAuthAPI("rbac", "GET", "/nav-menu/", GetRoleMenus)
}
