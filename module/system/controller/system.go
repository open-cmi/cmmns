package controller

import (
	"net/http"
	"os/exec"

	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/system/model"
	"github.com/open-cmi/goutils/devutil"

	"github.com/gin-gonic/gin"
)

// List list device info
func List(c *gin.Context) {
	var option api.Option
	api.ParseParams(c, &option)

	count, list, err := model.List(&option)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": list,
		}})
}

// Get get device info
func Get(c *gin.Context) {
	id := c.Param("id")
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	var option api.Option
	option.UserID = userID
	mdl := model.Get(&option, "id", id)
	if mdl == nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "not exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": mdl})
}

func Reboot(c *gin.Context) {

	exec.Command("/bin/sh", "-c", "reboot -f").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
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
