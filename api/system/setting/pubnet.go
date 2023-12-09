package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/setting/pubnet"
	"github.com/open-cmi/cmmns/service/webserver"
)

func SetPublicNet(c *gin.Context) {

	var req pubnet.SetPublicNetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": err.Error()})
		return
	}

	m := pubnet.Get()
	if m == nil {
		m = pubnet.New()
	}
	m.Host = req.Host
	m.Port = req.Port
	m.Schema = req.Schema
	err := m.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetPublicNet(c *gin.Context) {
	ex := pubnet.Get()
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": ex,
	})
}

func init() {
	webserver.RegisterAuthAPI("system-setting", "GET", "/pubnet/", GetPublicNet)
	webserver.RegisterAuthAPI("system-setting", "POST", "/pubnet/", SetPublicNet)
}
