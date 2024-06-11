package wac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/wac"
	"github.com/open-cmi/cmmns/module/wac/blacklist"
	"github.com/open-cmi/cmmns/module/wac/whitelist"
	"github.com/open-cmi/cmmns/pkg/goparam"
	"github.com/open-cmi/cmmns/service/webserver"
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
	var req wac.SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	if req.Enable {
		if req.Mode == "blacklist" {
			count, _, err := blacklist.List(nil)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
				return
			}
			if count == 0 {
				c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "blacklist must contains one address at least"})
				return
			}
		} else {
			count, _, err := whitelist.List(nil)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
				return
			}
			if count == 0 {
				c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "whitelist must contains one address at least"})
				return
			}
		}
	}
	err := wac.SetWAC(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func AddWhitelist(c *gin.Context) {
	var req wac.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := whitelist.AddWhitelist(req.Address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func DelWhitelist(c *gin.Context) {
	var req wac.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := whitelist.DelWhitelist(req.Address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func ListWhitelist(c *gin.Context) {
	param := goparam.ParseParams(c)

	count, lists, err := whitelist.List(param)
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
	var req wac.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := blacklist.AddBlacklist(req.Address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func DelBlacklist(c *gin.Context) {
	var req wac.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := blacklist.DelBlacklist(req.Address)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func ListBlacklist(c *gin.Context) {
	param := goparam.ParseParams(c)

	count, lists, err := blacklist.List(param)
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
	webserver.RegisterAuthRouter("wac", "/api/wac/v1/")
	webserver.RegisterAuthAPI("wac", "GET", "/", GetWAC)
	webserver.RegisterAuthAPI("wac", "POST", "/", SetWAC)

	webserver.RegisterAuthAPI("wac", "POST", "/blacklist/add/", AddBlacklist)
	webserver.RegisterAuthAPI("wac", "POST", "/blacklist/delete/", DelBlacklist)
	webserver.RegisterAuthAPI("wac", "GET", "/blacklist/", ListBlacklist)

	webserver.RegisterAuthAPI("wac", "POST", "/whitelist/add/", AddWhitelist)
	webserver.RegisterAuthAPI("wac", "POST", "/whitelist/delete/", DelWhitelist)
	webserver.RegisterAuthAPI("wac", "GET", "/whitelist/", ListWhitelist)
}
