package setting

import (
	"os/exec"

	"github.com/gin-gonic/gin"
)

func Reboot(c *gin.Context) {

	exec.Command("/bin/sh", "-c", "reboot -f").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func ShutDown(c *gin.Context) {
	exec.Command("/bin/sh", "-c", "shutdown -h now").Output()
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}
