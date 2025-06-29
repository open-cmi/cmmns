package tools

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/gobase/essential/webserver"
)

type MD5Request struct {
	Origin string `json:"origin"`
}

func MD5Encode(c *gin.Context) {
	var req MD5Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	result := fmt.Sprintf("%x", md5.Sum([]byte(req.Origin)))

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": result})
}

func init() {
	webserver.RegisterAuthAPI("tools", "POST", "/md5/encode/", MD5Encode)
}
