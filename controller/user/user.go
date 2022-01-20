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
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/auditlog"
	"github.com/open-cmi/cmmns/controller"

	model "github.com/open-cmi/cmmns/model/user"
	commsg "github.com/open-cmi/cmmns/msg/request"
	msg "github.com/open-cmi/cmmns/msg/user"
	"github.com/open-cmi/cmmns/utils"
	"github.com/open-cmi/goutils/typeutil"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/cmmns/storage/rdb"
)

// EmailTemplate html content template
var EmailTemplate string = `
<div>
	<h1>Hi username, Welcome to Nay!</h1>
	<h5>Here is a link to activate your account, please copy and paste it to your browser:</h5>
	<h5>https://domain/api/common/v3/user/activate/token</h5>
</div>
`

// CheckAuth get userinfo
func CheckAuth(c *gin.Context) {
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
	return
}

func ChangePassword(c *gin.Context) {
	var apimsg msg.ChangePasswordMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := controller.GetUser(c)
	if user == nil {
		return
	}

	if apimsg.NewPassword != apimsg.ConfirmPassword {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": "password confirmation doesn't match the password",
		})
		return
	}
	userID, _ := user["id"].(string)
	if !model.VerifyPasswordByID(userID, apimsg.OldPassword) {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": "user password verify failed",
		})
		return
	}

	err := model.ChangePassword(userID, apimsg.NewPassword)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": "change password failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

// List list user
func List(c *gin.Context) {

	var query commsg.RequestQuery
	utils.ParseParams(c, &query)

	count, users, err := model.List(&query)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": "list users failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": users,
		},
	})
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

	c.JSON(200, gin.H{"ret": 0, "msg": "", "data": user})
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

	cache := rdb.GetCache(rdb.UserCache)
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
	var apimsg msg.LoginMsg
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

	s, _ := c.Get("session")
	session, ok := s.(*sessions.Session)
	if ok {
		session.Values["user"] = map[string]interface{}{
			"username": user.UserName,
			"id":       user.ID,
			"email":    user.Email,
			"status":   user.Status,
			"role":     user.Role,
		}
	}

	// 写日志操作
	auditlog.InsertLog(c, auditlog.LoginType, "Login Success")

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": *user})

	return
}

func Logout(c *gin.Context) {
	sess, _ := c.Get("session")
	session := sess.(*sessions.Session)

	// 写日志操作
	auditlog.InsertLog(c, auditlog.LoginType, "Logout Success")

	delete(session.Values, "user")

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})

	return
}

// Register register user
func Register(c *gin.Context) {
	var apimsg msg.RegisterMsg
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
	if !typeutil.EmailIsValid(apimsg.Email) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email is not valid"})
		return
	}

	err := model.Register(&apimsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
	}

	code := uuid.New()
	cache := rdb.GetCache(rdb.UserCache)
	activateCode := fmt.Sprintf("activate_code_%s", code.String())
	err = cache.Set(context.TODO(), activateCode, apimsg.UserName, time.Hour*24).Err()
	if err != nil {
		model.DeleteByName(apimsg.UserName)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "code generate failed"})
		return
	}

	e := email.NewEmail()
	emailInfo := config.GetConfig().Email
	domain := config.GetConfig().MasterInfo.ExternalAddress
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
		model.DeleteByName(apimsg.UserName)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email can't be verified"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	return
}

// GetSelf get by self
func GetSelf(c *gin.Context) {
	cache, _ := c.Get("user")
	user, _ := cache.(model.BasicInfo)
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"username": user.UserName,
		},
	})
}

// CreateUser create user
func CreateUser(c *gin.Context) {
	var apimsg msg.CreateMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 验证邮箱有效性
	if !typeutil.EmailIsValid(apimsg.Email) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email is not valid"})
		return
	}

	err := model.Create(&apimsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	return
}

// DeleteUser delete user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := model.DeleteByID(id)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -1,
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"ret": 0, "msg": ""})
	return
}
