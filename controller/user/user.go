package user

import (
	"fmt"

	"github.com/open-cmi/cmmns/model/auth"

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

// Get get user by id
func Get(c *gin.Context) {

	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"username": "admin",
		},
	})
}

// GetSelf get by self
func GetSelf(c *gin.Context) {
	user, _ := c.Get("user")
	user = user.(auth.User)
	fmt.Println(user)
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"username": "admin",
		},
	})
}
