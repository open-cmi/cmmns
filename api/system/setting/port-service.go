package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/setting/service/web"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetServicePort(c *gin.Context) {
	m := web.GetServicePort()
	if m == nil {
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": *m})
}

func SetServicePort(c *gin.Context) {
	var req web.SetServicePortRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := web.SetServicePort(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthAPI("system", "POST", "/service/web-setting/", SetServicePort)
	webserver.RegisterAuthAPI("system", "GET", "/service/web-setting/", GetServicePort)
}
