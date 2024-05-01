package prod

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/system/prod"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetProdNav(c *gin.Context) {

	// todo, 管理员用户获取配置菜单
	// 普通用户获取授权菜单
	menu := prod.GetProdNav()

	// 返回LAN参数
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": menu,
	})
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/prod/nav/", GetProdNav)
}
