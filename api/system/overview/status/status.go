package status

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/system"
)

func GetStatus(c *gin.Context) {
	status := system.GetStatus()

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": status})
}
