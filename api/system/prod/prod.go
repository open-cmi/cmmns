package prod

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/setting/time"
	"github.com/open-cmi/cmmns/module/system/locale"
	"github.com/open-cmi/cmmns/module/system/prod"
	"github.com/open-cmi/gobase/essential/i18n"
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

func ResetProdBasisInfo(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	err := prod.ResetProdBasisInfo()
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("reset prod basis info"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("reset prod basis info"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetProdBasisInfo(c *gin.Context) {
	info, customized := prod.GetProdBasisInfo()
	m := time.Get()
	if m == nil {
		m = time.New()
	}

	_, hasCustomLogo := prod.GetLogoPath()
	_, hasCustomFavicon := prod.GetFaviconPath()
	_, hasCustomBg := prod.GetBackgroundPath()

	logoUpdatedAt := info.LogoUpdatedAt
	if logoUpdatedAt == 0 && hasCustomLogo {
		logoUpdatedAt = 1
	}
	favUpdatedAt := info.FaviconUpdatedAt
	if favUpdatedAt == 0 && hasCustomFavicon {
		favUpdatedAt = 1
	}
	bgUpdatedAt := info.BackgroundUpdatedAt
	if bgUpdatedAt == 0 && hasCustomBg {
		bgUpdatedAt = 1
	}

	data := gin.H{
		"name":               info.Name,
		"name_desc":          info.NameDesc,
		"footer":             info.Footer,
		"customized":         customized,
		"has_custom_logo":    hasCustomLogo,
		"logo_url":           fmt.Sprintf("/api/system/v1/prod/logo/?t=%d", logoUpdatedAt),
		"has_custom_favicon": hasCustomFavicon,
		"favicon_url":        fmt.Sprintf("/api/system/v1/prod/favicon/?t=%d", favUpdatedAt),
		"has_custom_bg":      hasCustomBg,
		"bg_url":             fmt.Sprintf("/api/system/v1/prod/background/?t=%d", bgUpdatedAt),
		"timezone":           m.TimeZone,
		"lang":               locale.GetLocale(),
	}

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": data,
	})
}

func init() {
	rbac.UnauthAPI("system", "GET", "/prod/basis/", GetProdBasisInfo)
	rbac.OptionAuthAPI("system", "POST", "/prod/basis/", SetProdBasisInfo, rbac.PermCode("cmmns:prodinfo", rbac.PermWrite))
	rbac.OptionAuthAPI("system", "DELETE", "/prod/basis/", ResetProdBasisInfo, rbac.PermCode("cmmns:prodinfo", rbac.PermWrite))
}
