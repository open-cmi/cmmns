package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/system/locale"
	"github.com/open-cmi/cmmns/service/webserver"
)

func ChangeLang(c *gin.Context) {
	type LangMsg struct {
		Lang string `json:"lang"`
	}

	var msg LangMsg
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	locale.SetLocale(msg.Lang)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetLang(c *gin.Context) {
	lang := locale.GetLocale()
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": lang,
	})
}

func init() {
	webserver.RegisterAuthAPI("system", "PUT", "/locale/", ChangeLang)
	webserver.RegisterAuthAPI("system", "GET", "/locale/", GetLang)
}
