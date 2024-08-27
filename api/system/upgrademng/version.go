package upgrademng

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/system/upgrademng"
	"github.com/open-cmi/cmmns/service/webserver"
)

func GetCurrentVersion(c *gin.Context) {
	v, err := upgrademng.CurrentVersion()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": v,
	})
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/upgrade-mng/current-version/", GetCurrentVersion)
}
