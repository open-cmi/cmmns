package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/webserver"
	"github.com/open-cmi/cmmns/pkg/dev"
)

func GetDevID(c *gin.Context) {

	deviceID := dev.GetDeviceID()

	// 返回LAN参数
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"id": deviceID,
		},
	})
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/device/", GetDevID)
}
