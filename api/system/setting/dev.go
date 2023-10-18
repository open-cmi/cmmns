package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/goutils/devutil"
)

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
