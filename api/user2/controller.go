package user2

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/open-cmi/cmmns/essential/rdb"
	"github.com/open-cmi/cmmns/essential/webserver"
	"github.com/open-cmi/cmmns/module/setting/email"
	"github.com/open-cmi/cmmns/module/user"
	"github.com/open-cmi/cmmns/pkg/verify"
)

// EmailTemplate html content template
var EmailTemplate string = `
<div>
	<h1>Hi username, Welcome to Nay!</h1>
	<h5>Here is a link to activate your account, please copy and paste it to your browser:</h5>
	<h5>user_activate_url/token</h5>
</div>
`

// Register register user
func Register(c *gin.Context) {
	var apimsg user.RegisterMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 验证验证码的有效性
	if !captcha.VerifyString(apimsg.CaptchaID, apimsg.Captcha) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "captcha is incorrect"})
		return
	}

	// 验证邮箱有效性
	if !verify.EmailIsValid(apimsg.Email) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email is not valid"})
		return
	}

	err := user.Register(&apimsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
	}

	code := uuid.New()
	cache := rdb.GetClient("user")
	activateCode := fmt.Sprintf("activate_code_%s", code.String())
	err = cache.Set(context.TODO(), activateCode, apimsg.UserName, time.Hour*24).Err()
	if err != nil {
		user.DeleteByName(apimsg.UserName)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "code generate failed"})
		return
	}

	activateURL := gConf.ActivateURL
	htmlcontent := strings.Replace(EmailTemplate, "token", code.String(), 1)
	htmlcontent = strings.Replace(htmlcontent, "user_activate_url", activateURL, 1)
	htmlcontent = strings.Replace(htmlcontent, "username", apimsg.UserName, 1)

	err = email.Send([]string{apimsg.Email}, []string{}, "Welcome to Nay", htmlcontent, nil)
	if err != nil {
		user.DeleteByName(apimsg.UserName)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email can't be verified"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// Activate activate user
func Activate(c *gin.Context) {
	code := c.Param("code")
	_, err := uuid.Parse(code)
	if err != nil {
		c.String(200, "activate code is not valid")
		return
	}

	cache := rdb.GetClient("user")
	activateCode := fmt.Sprintf("activate_code_%s", code)
	username, err := cache.Get(context.TODO(), activateCode).Result()
	if err != nil {
		c.String(200, "activate code is not exist")
		return
	}

	err = user.Activate(username)
	if err != nil {
		c.String(200, "activate user failed")
	} else {
		c.String(200, "activate user success, you can login now")
	}
}

// user v2版本适用于互联网用户，自己注册、激活等操作
func init() {
	webserver.RegisterAuthRouter("user2", "/api/user/v2")
	webserver.RegisterUnauthAPI("user2", "POST", "/register", Register)
	webserver.RegisterUnauthAPI("user2", "GET", "/activate/:code", Activate)
}
