package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/common/errcode"
	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/service/webserver"
	"github.com/open-cmi/cmmns/service/webserver/middleware"

	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/email"
	"github.com/open-cmi/cmmns/module/user"
	"github.com/open-cmi/goutils/typeutil"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/essential/rdb"
)

// EmailTemplate html content template
var EmailTemplate string = `
<div>
	<h1>Hi username, Welcome to Nay!</h1>
	<h5>Here is a link to activate your account, please copy and paste it to your browser:</h5>
	<h5>user_activate_url/token</h5>
</div>
`

// CheckAuth get userinfo
func CheckAuth(c *gin.Context) {
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func ChangePassword(c *gin.Context) {
	var apimsg user.ChangePasswordMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	usermap := goparam.GetUser(c)
	if usermap == nil {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("user not exist"),
		})
		return
	}

	if apimsg.NewPassword != apimsg.ConfirmPassword {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("password confirmation doesn't match the password"),
		})
		return
	}

	userID, _ := usermap["id"].(string)
	if !user.VerifyPasswordByID(userID, apimsg.OldPassword) {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("user password verify failed"),
		})
		return
	}

	err := user.ChangePassword(userID, apimsg.NewPassword)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("change password failed"),
		})
		return
	}

	ip := c.ClientIP()
	user := goparam.GetUser(c)
	if user != nil {
		username, _ := user["username"].(string)
		auditlog.InsertLog(ip,
			username,
			auditlog.OperationType,
			i18n.Sprintf("change password sussessfully"),
		)
	}

	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
	})
}

// List list user
func List(c *gin.Context) {

	var query goparam.Option
	goparam.ParseParams(c, &query)

	count, users, err := user.List(&query)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("list users failed"),
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

	id := c.Param("id")
	user := user.Get("id", id)
	if user == nil {
		c.JSON(200, gin.H{
			"ret": errcode.ErrFailed,
			"msg": i18n.Sprintf("user not exist"),
		})
		return
	}

	c.JSON(200, gin.H{"ret": 0, "msg": "", "data": user})
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

// Login login user
func Login(c *gin.Context) {
	var apimsg user.LoginMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 验证验证码的有效性
	if !apimsg.IgnoreCaptcha && !captcha.VerifyString(apimsg.CaptchaID, apimsg.Captcha) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "captcha is incorrect"})
		return
	}

	user, err := user.Login(&apimsg)
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
	ip := c.ClientIP()

	// 写日志操作
	auditlog.InsertLog(ip, user.UserName, auditlog.LoginType, i18n.Sprintf("login successfully"))

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": *user})
}

func Logout(c *gin.Context) {
	sess, _ := c.Get("session")
	session := sess.(*sessions.Session)

	ip := c.ClientIP()
	user := goparam.GetUser(c)
	if user != nil {
		username, _ := user["username"].(string)
		// 写日志操作
		auditlog.InsertLog(ip, username, auditlog.LoginType, i18n.Sprintf("logout successfully"))
	}

	delete(session.Values, "user")

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

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
	if !typeutil.EmailIsValid(apimsg.Email) {
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

	err = email.Send([]string{apimsg.Email}, "Welcome to Nay", htmlcontent, nil)
	if err != nil {
		user.DeleteByName(apimsg.UserName)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email can't be verified"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// Create create user
func Create(c *gin.Context) {
	var apimsg user.CreateMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 验证邮箱格式
	if !typeutil.EmailIsValid(apimsg.Email) {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email is not valid"})
		return
	}

	err := user.Create(&apimsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
	}

	ip := c.ClientIP()
	user := goparam.GetUser(c)
	if user != nil {
		username, _ := user["username"].(string)
		auditlog.InsertLog(ip,
			username,
			auditlog.OperationType,
			i18n.Sprintf("create user %s sussessfully", apimsg.UserName),
		)
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// Delete delete user
func Delete(c *gin.Context) {
	id := c.Param("id")
	err := user.DeleteByID(id)
	if err != nil {
		c.JSON(200, gin.H{
			"ret": -1,
			"msg": err.Error(),
		})
		return
	}
	ip := c.ClientIP()
	user := goparam.GetUser(c)
	if user != nil {
		username, _ := user["username"].(string)
		// 写日志操作
		auditlog.InsertLog(ip,
			username,
			auditlog.OperationType,
			i18n.Sprintf("delete user sussessfully"),
		)
	}

	c.JSON(200, gin.H{"ret": 0, "msg": ""})
}

func CreateToken(c *gin.Context) {
	user := goparam.GetUser(c)
	if user == nil {
		c.JSON(200, gin.H{"ret": -1, "msg": "user data is empty"})
		return
	}

	username, _ := user["username"].(string)
	userid, _ := user["id"].(string)
	email, _ := user["email"].(string)
	role, _ := user["role"].(int)
	status, _ := user["status"].(int)
	tk, err := middleware.GenerateAuthToken(username, userid, email, role, status, 30)
	if err != nil {
		c.JSON(200, gin.H{"ret": -1, "msg": "create token failed"})
		return
	}
	c.JSON(200, gin.H{"ret": 0, "msg": "", "token": tk})
}

func init() {
	webserver.RegisterAuthRouter("user", "/api/common/v3/user")
	webserver.RegisterAuthAPI("user", "GET", "/checkauth", CheckAuth)
	webserver.RegisterAuthAPI("user", "GET", "/", List)
	webserver.RegisterAuthAPI("user", "POST", "/", Create)
	webserver.RegisterAuthAPI("user", "POST", "/changepass", ChangePassword)
	webserver.RegisterAuthAPI("user", "POST", "/logout", Logout)
	webserver.RegisterAuthAPI("user", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("user", "DELETE", "/:id", Delete)

	webserver.RegisterUnauthRouter("user", "/api/common/v3/user")
	webserver.RegisterUnauthAPI("user", "POST", "/login", Login)
	webserver.RegisterUnauthAPI("user", "POST", "/register", Register)
	webserver.RegisterUnauthAPI("user", "GET", "/activate/:code", Activate)
	webserver.RegisterUnauthAPI("user", "POST", "/jwt-token/", CreateToken)
}
