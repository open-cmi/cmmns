package user

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/dchest/captcha"
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/module/middleware"
	"github.com/open-cmi/cmmns/module/setting/pubnet"
	"github.com/open-cmi/cmmns/pkg/goparam"
	"github.com/open-cmi/cmmns/pkg/verify"
	"github.com/open-cmi/cmmns/service/webserver"

	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/user"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/essential/logger"
)

// CheckAuth get userinfo
func CheckAuth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

var PasswordComplexRegexp *regexp.Regexp

func VerifyPasswordRule(str string, minLen, maxLen int) error {
	var (
		isUpper   = false
		isLower   = false
		isNumber  = false
		isSpecial = false
	)

	if len(str) < minLen || len(str) > maxLen {
		return errors.New("the password must contain uppercase and lowercase letters, numbers or punctuation, and must be 6-30 digits long. ")
	}

	for _, s := range str {
		switch {
		case unicode.IsUpper(s):
			isUpper = true
		case unicode.IsLower(s):
			isLower = true
		case unicode.IsNumber(s):
			isNumber = true
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			isSpecial = true
		default:
		}
	}

	if (isUpper && isLower) && (isNumber || isSpecial) {
		return nil
	}
	return errors.New("the password must contain uppercase and lowercase letters, numbers or punctuation, and must be 6-30 digits long. ")
}

func ChangePassword(c *gin.Context) {
	var apimsg user.ChangePasswordMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	usermap := goparam.GetUser(c)
	if usermap == nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("user not exist"),
		})
		return
	}

	if apimsg.NewPassword != apimsg.ConfirmPassword {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("password confirmation doesn't match the password"),
		})
		return
	}

	userID, _ := usermap["id"].(string)
	if !user.VerifyPasswordByID(userID, apimsg.OldPassword) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("user password verify failed"),
		})
		return
	}
	err := VerifyPasswordRule(apimsg.NewPassword, 8, 30)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf(err.Error()),
		})
		return
	}
	err = user.ChangePassword(userID, apimsg.NewPassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("change password failed"),
		})
		return
	}

	ip := c.ClientIP()
	username, _ := usermap["username"].(string)
	auditlog.InsertLog(ip,
		username,
		auditlog.OperationType,
		i18n.Sprintf("change password sussessfully"),
	)

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

// List list user
func List(c *gin.Context) {

	query := goparam.ParseParams(c)

	count, users, err := user.List(query)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("list users failed"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
		c.JSON(http.StatusOK, gin.H{
			"ret": -1,
			"msg": i18n.Sprintf("user not exist"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": user})
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

	// 记录公网ip
	m := pubnet.Get()
	if m == nil {
		m = pubnet.New()
		host := c.Request.Host
		hostport := strings.Split(host, ":")
		m.Host = hostport[0]
		if hostport[1] != "" {
			m.Port, _ = strconv.Atoi(hostport[1])
		}

		if c.Request.TLS != nil {
			m.Schema = "https"
		} else if proto := c.GetHeader("X-Forwarded-Proto"); proto == "https" {
			m.Schema = "https"
		} else {
			m.Schema = "http"
		}

		err = m.Save()
		if err != nil {
			logger.Errorf("save pubnet failed: %s\n", err.Error())
		}
	}
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

	session.Options.MaxAge = -1 // aged
	delete(session.Values, "user")

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
	if !verify.EmailIsValid(apimsg.Email) {
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
	u := goparam.GetUser(c)
	id := c.Param("id")
	userID := u["id"].(string)
	if id == userID {
		c.JSON(http.StatusOK, gin.H{
			"ret": -1,
			"msg": "can't delete youself",
		})
		return
	}

	err := user.Delete(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": -1,
			"msg": err.Error(),
		})
		return
	}

	ip := c.ClientIP()
	if u != nil {
		username, _ := u["username"].(string)
		// 写日志操作
		auditlog.InsertLog(ip,
			username,
			auditlog.OperationType,
			i18n.Sprintf("delete user sussessfully"),
		)
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func CreateToken(c *gin.Context) {
	var req middleware.CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := goparam.GetUser(c)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "user data is empty"})
		return
	}

	username, _ := user["username"].(string)
	userid, _ := user["id"].(string)
	email, _ := user["email"].(string)
	role, _ := user["role"].(int)
	status, _ := user["status"].(int)
	tk, err := middleware.GenerateAuthToken(req.Name, username, userid, email, role, status, req.ExpireDay)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "create token failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "token": tk})
}

func Edit(c *gin.Context) {
	var req user.EditMsg
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := user.Edit(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func ResetPassword(c *gin.Context) {
	var req user.ResetPasswdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	usermap := goparam.GetUser(c)
	if usermap == nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("user not exist"),
		})
		return
	}

	if req.Password != req.Password2 {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("password confirmation doesn't match the password"),
		})
		return
	}

	err := user.ResetPasswd(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func TokenList(c *gin.Context) {
	query := goparam.ParseParams(c)

	count, tokens, err := middleware.TokenList(query)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("list tokens failed"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": tokens,
		},
	})
}

// v1 user适用于普通的后台管理系统，用户由管理员创建管理，不支持自注册用户
func init() {
	webserver.RegisterAuthRouter("user", "/api/user/v1")
	webserver.RegisterAuthAPI("user", "GET", "/checkauth", CheckAuth)
	webserver.RegisterAuthAPI("user", "GET", "/", List)
	webserver.RegisterAuthAPI("user", "POST", "/", Create)
	webserver.RegisterAuthAPI("user", "POST", "/change-passwd", ChangePassword)
	webserver.RegisterAuthAPI("user", "POST", "/reset-passwd", ResetPassword)
	webserver.RegisterAuthAPI("user", "POST", "/logout", Logout)
	webserver.RegisterAuthAPI("user", "GET", "/:id", Get)
	webserver.RegisterAuthAPI("user", "PUT", "/:id", Edit)
	webserver.RegisterAuthAPI("user", "DELETE", "/:id", Delete)
	webserver.RegisterAuthAPI("user", "POST", "/jwt-token/", CreateToken)
	webserver.RegisterAuthAPI("user", "GET", "/jwt-token/", TokenList)

	webserver.RegisterUnauthRouter("user", "/api/user/v1")
	webserver.RegisterUnauthAPI("user", "POST", "/login", Login)
}
