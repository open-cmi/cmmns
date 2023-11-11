package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/module/setting"
	"github.com/open-cmi/cmmns/module/setting/pubnet"
)

func List(c *gin.Context) {
	var option goparam.Option
	goparam.ParseParams(c, &option)

	count, results, err := setting.List(&option)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": results,
		},
	})
}

func Get(c *gin.Context) {
	identify := c.Param("id")

	user := goparam.GetUser(c)
	userID, _ := user["id"].(string)
	var option goparam.Option
	option.UserID = userID

	m := setting.FilterGet(&option, []string{"id"}, []interface{}{identify})

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"date": m,
	})
}

func Edit(c *gin.Context) {
	identify := c.Param("id")

	var reqMsg setting.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := goparam.GetUser(c)
	userID, _ := user["id"].(string)
	var option goparam.Option
	option.UserID = userID

	err := setting.Edit(&option, identify, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func SetPublicNet(c *gin.Context) {

	var req pubnet.SetPublicNetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": err.Error()})
		return
	}

	m := pubnet.Get()
	if m == nil {
		m = pubnet.New()
	}
	m.Host = req.Host
	m.Port = req.Port
	m.Schema = req.Schema
	err := m.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func GetPublicNet(c *gin.Context) {
	ex := pubnet.Get()
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": ex,
	})
}
