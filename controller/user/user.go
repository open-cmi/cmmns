package user

import (
	"context"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/google/uuid"
	climsg "github.com/open-cmi/cmmns/climsg/user"
	model "github.com/open-cmi/cmmns/model/user"
	"github.com/open-cmi/goutils/verify"

	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/cmmns/db"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
)

// EmailTemplate html content template
var EmailTemplate string = `
<div>
	<h1>Hi username, Welcome to Nay!</h1>
	<h5>Here is a link to activate your account, please copy and paste it to your browser:</h5>
	<h5>https://domain/api/common/v3/user/activate/token</h5>
</div>
`

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

	expire := time.Now().Add(10 * time.Minute)
	cookie := http.Cookie{
		Name:     "test",
		Value:    "this is a test",
		Expires:  expire,
		Path:     "/",
		HttpOnly: false,
	}

	id := c.Param("id")
	user, err := model.Get(id)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -1,
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"ret": 0, "msg": "", "data": *user})
	http.SetCookie(c.Writer, &cookie)
	return
}

// Activate activate user
func Activate(c *gin.Context) {
	code := c.Param("code")
	fmt.Println("code:", code)
	_, err := uuid.Parse(code)
	if err != nil {
		c.String(200, "activate code is not valid")
		return
	}

	cache := db.GetCache()
	activateCode := fmt.Sprintf("activate_code_%s", code)
	username, err := cache.Get(context.TODO(), activateCode).Result()
	if err != nil {
		c.String(200, "activate code is not exist")
		return
	}

	err = model.Activate(username)
	if err != nil {
		c.String(200, "activate user failed")
	} else {
		c.String(200, "activate user success, you can login now")
	}
	return
}

// Login login user
func Login(c *gin.Context) {
	var apimsg climsg.LoginMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 验证验证码的有效性
	if !apimsg.IgnoreCaptcha && !captcha.VerifyString(apimsg.CaptchaID, apimsg.Captcha) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "captcha is incorrect"})
		return
	}

	user, err := model.Login(&apimsg)
	if err != nil {
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

	// 验证验证码的有效性
	if !apimsg.IgnoreCaptcha && !captcha.VerifyString(apimsg.CaptchaID, apimsg.Captcha) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "captcha is incorrect"})
		return
	}

	// 验证邮箱有效性
	if !verify.EmailIsValid(apimsg.Email) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email is not valid"})
		return
	}

	err := model.Register(&apimsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
	}

	code := uuid.New()
	cache := db.GetCache()
	activateCode := fmt.Sprintf("activate_code_%s", code.String())
	err = cache.Set(context.TODO(), activateCode, apimsg.UserName, time.Hour*24).Err()
	if err != nil {
		model.Delete(apimsg.UserName)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "code generate failed"})
		return
	}

	e := email.NewEmail()
	emailInfo := config.GetConfig().Email
	domain := config.GetConfig().Domain
	e.From = emailInfo.From
	e.To = []string{apimsg.Email}
	//e.Cc = []string{"danielzhao2012@163.com"}
	e.Subject = "Welcome to Nay"
	//e.Text = []byte("Text Body is, of course, supported!")
	htmlcontent := strings.Replace(EmailTemplate, "token", code.String(), 1)
	htmlcontent = strings.Replace(htmlcontent, "domain", domain, 1)
	htmlcontent = strings.Replace(htmlcontent, "username", apimsg.UserName, 1)

	e.HTML = []byte(htmlcontent)
	err = e.Send(emailInfo.SMTPServer, smtp.PlainAuth("", emailInfo.User, emailInfo.Password, emailInfo.SMTPHost))
	if err != nil {
		model.Delete(apimsg.UserName)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email can't be verified"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	return
}

// GetSelf get by self
func GetSelf(c *gin.Context) {
	cache, _ := c.Get("user")
	user, _ := cache.(model.User)
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"username": user.UserName,
		},
	})
}
