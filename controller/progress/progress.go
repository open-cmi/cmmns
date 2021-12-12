package progress

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get get progress
func Get(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"total":   0,
			"success": 0,
			"failed":  0,
		}})
	return
}
