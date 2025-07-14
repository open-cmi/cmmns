package wac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/wac"
	"github.com/open-cmi/cmmns/module/wac/blacklist"
	"github.com/open-cmi/cmmns/module/wac/whitelist"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/webserver"
	"github.com/open-cmi/gobase/pkg/goparam"
)

func GetWAC(c *gin.Context) {
	m := wac.GetWAC()
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": m,
	})
}

func SetWAC(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	var req wac.SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set web access control"), false)
		return
	}
	if req.Enable {
		if req.Mode == "blacklist" {
			count, _, err := blacklist.QueryList(nil)
			if err != nil {
				ah.InsertOperationLog(i18n.Sprintf("set web access control"), false)
				c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
				return
			}
			if count == 0 {
				ah.InsertOperationLog(i18n.Sprintf("set web access control"), false)
				c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "blacklist must contains one address at least"})
				return
			}
		} else {
			count, _, err := whitelist.QueryList(nil)
			if err != nil {
				ah.InsertOperationLog(i18n.Sprintf("set web access control"), false)
				c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
				return
			}
			if count == 0 {
				ah.InsertOperationLog(i18n.Sprintf("set web access control"), false)
				c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "whitelist must contains one address at least"})
				return
			}
		}
	}
	err := wac.SetWAC(&req)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("set web access control"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("set web access control"), true)

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func AddWhitelist(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req wac.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("web access control add whitelist"), false)
		return
	}
	err := whitelist.AddWhitelist(req.Address)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("web access control add whitelist"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("web access control add whitelist"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func DelWhitelist(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req wac.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("web access control delete whitelist"), false)
		return
	}
	err := whitelist.DelWhitelist(req.Address)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("web access control delete whitelist"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("web access control delete whitelist"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func ListWhitelist(c *gin.Context) {
	param := goparam.ParseParams(c)

	count, lists, err := whitelist.QueryList(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": lists,
		},
	})
}

func AddBlacklist(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req wac.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("web access control add blacklist"), false)
		return
	}
	err := blacklist.AddBlacklist(req.Address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("web access control add blacklist"), false)
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("web access control add blacklist"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func DelBlacklist(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req wac.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("web access control delete blacklist"), false)
		return
	}

	err := blacklist.DelBlacklist(req.Address)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("web access control delete blacklist"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("web access control delete blacklist"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func ListBlacklist(c *gin.Context) {
	param := goparam.ParseParams(c)

	count, lists, err := blacklist.QueryList(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": lists,
		},
	})
}

func init() {
	webserver.RegisterMustAuthRouter("wac", "/api/wac/v1/")
	webserver.RegisterMustAuthAPI("wac", "GET", "/", GetWAC)
	webserver.RegisterMustAuthAPI("wac", "POST", "/", SetWAC)

	webserver.RegisterMustAuthAPI("wac", "POST", "/blacklist/add/", AddBlacklist)
	webserver.RegisterMustAuthAPI("wac", "POST", "/blacklist/delete/", DelBlacklist)
	webserver.RegisterMustAuthAPI("wac", "GET", "/blacklist/", ListBlacklist)

	webserver.RegisterMustAuthAPI("wac", "POST", "/whitelist/add/", AddWhitelist)
	webserver.RegisterMustAuthAPI("wac", "POST", "/whitelist/delete/", DelWhitelist)
	webserver.RegisterMustAuthAPI("wac", "GET", "/whitelist/", ListWhitelist)
}
