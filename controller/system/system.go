package system

import (
	"net/http"
	"os/exec"

	"github.com/open-cmi/cmmns/controller"
	model "github.com/open-cmi/cmmns/model/system"

	"github.com/gin-gonic/gin"
)

// List list device info
func List(c *gin.Context) {
	var option model.Option
	controller.ParseParams(c, &option.Option)

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
	return
}

// Get get device info
func Get(c *gin.Context) {
	id := c.Param("id")
	user := controller.GetUser(c)
	userID, _ := user["id"].(string)

	var option model.Option
	option.Option.UserID = userID
	mdl := model.Get(&option, "id", id)
	if mdl == nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "not exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": mdl})
	return
}

func Reboot(c *gin.Context) {

	exec.Command("/bin/sh", "-c", "reboot -f").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}
