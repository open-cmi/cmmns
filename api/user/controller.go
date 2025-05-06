package user

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/dchest/captcha"
	"github.com/gorilla/sessions"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/essential/webserver"
	"github.com/open-cmi/cmmns/module/middleware"
	"github.com/open-cmi/cmmns/module/setting/pubnet"
	"github.com/open-cmi/cmmns/pkg/goparam"
	"github.com/open-cmi/cmmns/pkg/verify"

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
	ah := auditlog.NewAuditHandler(c)
	var apimsg user.ChangePasswordMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	usermap := goparam.GetUser(c)
	if usermap == nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("user does not exist"),
		})
		ah.InsertOperationLog(i18n.Sprintf("change password"), false)
		return
	}

	if apimsg.NewPassword != apimsg.ConfirmPassword {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("password confirmation doesn't match the password"),
		})

		ah.InsertOperationLog(i18n.Sprintf("change password"), false)
		return
	}

	userID, _ := usermap["id"].(string)
	if !user.VerifyPasswordByID(userID, apimsg.OldPassword) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("user password verify failed"),
		})

		ah.InsertOperationLog(i18n.Sprintf("change password"), false)
		return
	}
	err := VerifyPasswordRule(apimsg.NewPassword, 8, 30)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf(err.Error()),
		})

		ah.InsertOperationLog(i18n.Sprintf("change password"), false)
		return
	}
	err = user.ChangePassword(userID, apimsg.NewPassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("change password failed"),
		})

		ah.InsertOperationLog(i18n.Sprintf("change password"), false)
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("change password"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

// List list user
func List(c *gin.Context) {
	query := goparam.ParseParams(c)
	var paramnum int = 1
	var whereClause string
	var whereArgs []interface{}

	username := c.Query("username")
	if username != "" {
		if whereClause != "" {
			whereClause += " and "
		}
		whereClause += fmt.Sprintf(`username like %s`, sqldb.LikePlaceHolder(paramnum))
		whereArgs = append(whereArgs, username)
		paramnum += 1
	}

	if paramnum != 1 {
		query.WhereClause = " where " + whereClause
		query.WhereArgs = whereArgs
	}

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
	user := user.Get(id)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": -1,
			"msg": i18n.Sprintf("user does not exist"),
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
	ah := auditlog.NewAuditHandler(c)

	user, err := user.Login(&apimsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertLoginLog(apimsg.UserName, i18n.Sprintf("login"), false)
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
	ah.InsertLoginLog(user.UserName, i18n.Sprintf("login"), true)

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

	ah := auditlog.NewAuditHandler(c)
	ah.InsertOperationLog(i18n.Sprintf("logout"), true)

	session.Options.MaxAge = -1 // aged
	delete(session.Values, "user")

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// Create create user
func Create(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var apimsg user.CreateMsg
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("create user"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 验证邮箱格式
	if !verify.EmailIsValid(apimsg.Email) {
		ah.InsertOperationLog(i18n.Sprintf("create user %s", apimsg.UserName), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "email is not valid"})
		return
	}

	err := user.Create(&apimsg)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("create user %s", apimsg.UserName), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("create user %s", apimsg.UserName), true)

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// Delete delete user
func Delete(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	param := goparam.ParseParams(c)
	id := c.Param("id")
	userID := param.UserID
	if id == userID {
		c.JSON(http.StatusOK, gin.H{
			"ret": -1,
			"msg": "can't delete youself",
		})
		ah.InsertOperationLog(i18n.Sprintf("delete user"), false)
		return
	}
	role := param.Role
	if role != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"ret": -1,
			"msg": i18n.Sprintf("no permission"),
		})
		ah.InsertOperationLog(i18n.Sprintf("delete user"), false)
		return
	}

	err := user.Delete(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": -1,
			"msg": err.Error(),
		})
		ah.InsertOperationLog(i18n.Sprintf("delete user"), false)
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("delete user"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func CreateToken(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	var req middleware.CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("create token"), false)
		return
	}

	user := goparam.GetUser(c)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "user data is empty"})
		ah.InsertOperationLog(i18n.Sprintf("create token"), false)
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
		ah.InsertOperationLog(i18n.Sprintf("create token"), false)
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("create token"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "token": tk})
}

func DeleteToken(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req middleware.DeleteTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("delete token"), false)
		return
	}

	err := middleware.DeleteAuthToken(req.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "delete token failed"})
		ah.InsertOperationLog(i18n.Sprintf("delete token"), false)
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("delete token"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func Edit(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req user.EditMsg
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("edit user"), false)
		return
	}
	err := user.Edit(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("edit user %s", req.Username), false)
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("edit user %s", req.Username), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func ResetPassword(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req user.ResetPasswdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("reset password"), false)
		return
	}
	usermap := goparam.GetUser(c)
	if usermap == nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("user does not exist"),
		})
		ah.InsertOperationLog(i18n.Sprintf("reset password"), false)
		return
	}

	if req.Password != req.Password2 {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("password confirmation doesn't match the password"),
		})
		ah.InsertOperationLog(i18n.Sprintf("reset password"), false)
		return
	}

	err := user.ResetPasswd(&req)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("reset password"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("reset password"), true)
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
	webserver.RegisterMustAuthRouter("user", "/api/user/v1")
	webserver.RegisterMustAuthAPI("user", "GET", "/checkauth", CheckAuth)
	webserver.RegisterMustAuthAPI("user", "GET", "/", List)
	webserver.RegisterMustAuthAPI("user", "POST", "/", Create)
	webserver.RegisterMustAuthAPI("user", "POST", "/change-passwd", ChangePassword)
	webserver.RegisterMustAuthAPI("user", "POST", "/reset-passwd", ResetPassword)
	webserver.RegisterMustAuthAPI("user", "POST", "/logout", Logout)
	webserver.RegisterMustAuthAPI("user", "GET", "/:id", Get)
	webserver.RegisterMustAuthAPI("user", "PUT", "/:id", Edit)
	webserver.RegisterMustAuthAPI("user", "DELETE", "/:id", Delete)
	webserver.RegisterMustAuthAPI("user", "POST", "/jwt-token/", CreateToken)
	webserver.RegisterMustAuthAPI("user", "POST", "/jwt-token/delete/", DeleteToken)
	webserver.RegisterMustAuthAPI("user", "GET", "/jwt-token/", TokenList)

	webserver.RegisterUnauthRouter("user", "/api/user/v1")
	webserver.RegisterUnauthAPI("user", "POST", "/login", Login)
}
