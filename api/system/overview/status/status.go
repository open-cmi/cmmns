package status

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/system"
	"github.com/open-cmi/goutils/devutil"
)

func GetStatus(c *gin.Context) {
	status := system.GetStatus()

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": status})
}

func GetDevID(c *gin.Context) {

	deviceID := devutil.GetDeviceID()

	// 返回LAN参数
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"id": deviceID,
		},
	})
}
