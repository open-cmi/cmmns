package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/network"
	"github.com/open-cmi/cmmns/service/webserver"
)

func SetNetwork(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var conf network.ConfigRequest

	if err := c.ShouldBindJSON(&conf); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set network"), false)
		return
	}

	err := network.Set(&conf)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("set network"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("set network"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetNetwork(c *gin.Context) {

	conf := network.Get()

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": map[string]interface{}{
		"setting": conf,
	}})
}

func GetNetworkStatus(c *gin.Context) {
	count, status, err := network.GetStatus()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": status,
		}})
}

func BlinkingNetworkInterface(c *gin.Context) {
	var req network.BlinkingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	err := network.BlinkingInterface(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func SetManagementInterface(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req network.SetManagementInterfaceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set management interface"), false)
		return
	}

	err := network.SetManagementInterface(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set management interface"), false)

		return
	}
	ah.InsertOperationLog(i18n.Sprintf("set management interface"), true)

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetAvailableManagementInterface(c *gin.Context) {
	devices, err := network.GetAvailableManagementInterface()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": devices,
	})
}

func init() {
	webserver.RegisterAuthRouter("network", "/api/network/v1/")
	webserver.RegisterAuthAPI("network", "GET", "/", GetNetwork)
	webserver.RegisterAuthAPI("network", "POST", "/", SetNetwork)
	webserver.RegisterAuthAPI("network", "GET", "/status/", GetNetworkStatus)
	webserver.RegisterAuthAPI("network", "POST", "/blinking/", BlinkingNetworkInterface)
	webserver.RegisterAuthAPI("network", "POST", "/management-interface/", SetManagementInterface)
	webserver.RegisterAuthAPI("network", "GET", "/available-management-interface/", GetAvailableManagementInterface)
}
