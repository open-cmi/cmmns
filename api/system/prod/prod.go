package prod

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/setting/time"
	"github.com/open-cmi/cmmns/module/system/locale"
	"github.com/open-cmi/cmmns/module/system/prod"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetProdInfo(c *gin.Context) {

	// todo, 管理员用户获取配置菜单
	// 普通用户获取授权菜单
	info := prod.GetProdInfo()
	m := time.Get()

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]string{
			"name":     info.Name,
			"footer":   info.Footer,
			"timezone": m.TimeZone,
			"lang":     locale.GetLocale(),
		},
	})
}

func init() {
	webserver.RegisterUnauthAPI("system", "GET", "/prod/brief-info/", GetProdInfo)
}
