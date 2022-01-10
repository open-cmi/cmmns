package agent

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/cmmns/controller/ctl"
)

type MasterSetting struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
}

type Setting struct {
	Master MasterSetting `json:"master"`
}

// 获取配置
func GetSetting(c *gin.Context) {
	if config.GetConfig().MasterInfo.ExternalAddress == "" {
		var address string = ""
		host := c.Request.Host
		// 目前只支持ipv4地址，不支持ipv6地址
		if strings.Contains(host, ":") {
			arr := strings.Split(host, ":")
			address = arr[0]
		} else {
			address = host
		}
		config.GetConfig().MasterInfo.ExternalAddress = address
		config.GetConfig().Save()
	}
	var setting Setting
	setting.Master.Address = config.GetConfig().MasterInfo.ExternalAddress
	setting.Master.Port = config.GetConfig().MasterInfo.ExternalPort
	setting.Master.Proto = config.GetConfig().MasterInfo.ExternalProto

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": setting,
	})
}

func EditSetting(c *gin.Context) {
	user := ctl.GetUser(c)
	role, _ := user["role"].(int)
	if role != 0 {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": "no permision",
		})
		return
	}

	var reqMsg Setting
	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	config.GetConfig().MasterInfo.ExternalAddress = reqMsg.Master.Address
	config.GetConfig().MasterInfo.ExternalPort = reqMsg.Master.Port
	config.GetConfig().MasterInfo.ExternalProto = reqMsg.Master.Proto

	config.GetConfig().Save()

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
	return
}
