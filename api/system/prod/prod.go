package prod

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/setting/time"
	"github.com/open-cmi/cmmns/module/system/locale"
	"github.com/open-cmi/cmmns/module/system/prod"
	"github.com/open-cmi/cmmns/service/webserver"
)

func SetProdBasisInfo(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	var req prod.ProdInfoSetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("set prod basis info"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	err := prod.SetProdBasisInfo(&req)
	if err != nil {

		ah.InsertOperationLog(i18n.Sprintf("set prod basis info"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("set prod basis info"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetProdBasisInfo(c *gin.Context) {

	// todo, 管理员用户获取配置菜单
	// 普通用户获取授权菜单
	info := prod.GetProdBasisInfo()
	m := time.Get()
	if m == nil {
		m = time.New()
	}

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

func ToggleExperimentalSetting(c *gin.Context) {

	ah := auditlog.NewAuditHandler(c)
	ah.InsertOperationLog(i18n.Sprintf("toggle experimental setting"), true)
	prod.ToggleExperimentalSetting()

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func init() {
	webserver.RegisterUnauthAPI("system", "GET", "/prod/basis/", GetProdBasisInfo)
	webserver.RegisterUnauthAPI("system", "POST", "/prod/basis/", SetProdBasisInfo)
	webserver.RegisterAuthAPI("system", "POST", "/prod/experimental-setting/", ToggleExperimentalSetting)
}
