package netconfig

import (
	"github.com/gin-gonic/gin"
)

// List list user
func List(c *gin.Context) {
	c.JSON(200, gin.H{
		"ret":  0,
		"msg":  "",
		"data": []interface{}{},
	})
}

// GetEth get eth by name
func GetEth(c *gin.Context) {
	eth := c.Param("eth")
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"username": "admin",
		},
	})
}

// Save save all eth
func Save(c *gin.Context) {
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"username": "admin",
		},
	})
}

// SaveEth save eth
func SaveEth(c *gin.Context) {
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"username": "admin",
		},
	})
}
