package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/system/locale"
	"github.com/open-cmi/gobase/essential/i18n"
)

func ChangeLang(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	type LangMsg struct {
		Lang string `json:"lang"`
	}

	var msg LangMsg

	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set language"), false)
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("set language"), true)

	locale.SetLocale(msg.Lang)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetLang(c *gin.Context) {
	lang := locale.GetLocale()
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": lang,
	})
}

func init() {
	rbac.OptionAuthAPI("system", "PUT", "/locale/", ChangeLang, rbac.GetInitRoles())
	rbac.OptionAuthAPI("system", "GET", "/locale/", GetLang, rbac.GetInitRoles())
}
