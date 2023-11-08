package wac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/wac"
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
	err := wac.SetWAC(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthRouter("wac", "/api/wac/v1/")
	webserver.RegisterAuthAPI("wac", "GET", "/", GetWAC)
	webserver.RegisterAuthAPI("wac", "POST", "/", SetWAC)
}
