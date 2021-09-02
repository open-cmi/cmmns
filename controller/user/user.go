package user

import (
	"fmt"
	"net/http"

	"github.com/dchest/captcha"
	climsg "github.com/open-cmi/cmmns/climsg/user"
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

// Login login user
func Login(c *gin.Context) {
	var apimsg climsg.LoginMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	fmt.Println(apimsg)
	// 验证验证码的有效性
	if !apimsg.IgnoreCaptcha && !captcha.VerifyString(apimsg.CaptchaID, apimsg.Captcha) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "captcha is incorrect"})
		return
	}

	user, err := model.Login(&apimsg)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.Set("user", user)

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": user})
	return
}

// Register register user
func Register(c *gin.Context) {
	var apimsg climsg.RegisterMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	fmt.Println(apimsg)
	// 验证验证码的有效性
	if !apimsg.IgnoreCaptcha && !captcha.VerifyString(apimsg.CaptchaID, apimsg.Captcha) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "captcha is incorrect"})
		return
	}

	err := model.Register(&apimsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	}
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
