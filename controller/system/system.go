package system

import (
	"os/exec"

	model "github.com/open-cmi/cmmns/model/system"

	"github.com/gin-gonic/gin"
)

// GetSystemInfo get device info
func GetSystemInfo(c *gin.Context) {
	info, err := model.GetSystemInfo()
	if err != nil {
		c.JSON(200, gin.H{"ret": 1, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"ret":  0,
		"msg":  "",
		"data": info,
	})
}

func Reboot(c *gin.Context) {

	exec.Command("/bin/sh", "-c", "reboot").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}
