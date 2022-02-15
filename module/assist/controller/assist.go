package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/assist/infra"
)

func GetAssist(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": infra.IsRunning(),
	})
}

func SetAssist(c *gin.Context) {

	var msg struct {
		Enable bool `json:"enable"`
	}

	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	if !msg.Enable {
		infra.Close()
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "",
		})
		return
	}

	err := infra.Run()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": "start assist service failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
