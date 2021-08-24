package user

import (
	"fmt"

	"github.com/open-cmi/cmmns/model/auth"
	model "github.com/open-cmi/cmmns/model/user"

	"github.com/gin-gonic/gin"
)

// List list user
func List(c *gin.Context) {
	users, err := model.List()
	if err != nil {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": "list users failed",
		})
	} else {
		c.JSON(200, gin.H{
			"ret":  0,
			"msg":  "",
			"data": users,
		})
	}
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

// Register register user
func Register(c *gin.Context) {

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
