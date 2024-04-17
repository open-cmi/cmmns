package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/network"
)

func SetNetwork(c *gin.Context) {
	var conf network.ConfigRequest

	if err := c.ShouldBindJSON(&conf); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	err := network.Set(&conf)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetNetwork(c *gin.Context) {

	conf := network.Get()

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": map[string]interface{}{
		"setting": conf,
	}})
}

func GetNetworkStatus(c *gin.Context) {
	status, err := network.GetStatus()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": map[string]interface{}{
		"status": status,
	}})
}
