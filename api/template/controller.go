package template

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/parameter"
	"github.com/open-cmi/cmmns/module/template"
)

func List(c *gin.Context) {
	var option parameter.Option
	parameter.ParseParams(c, &option)

	count, results, err := template.List(&option)
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

func MultiDelete(c *gin.Context) {
	var reqMsg template.MultiDeleteMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := parameter.GetUser(c)
	userID, _ := user["id"].(string)

	var option parameter.Option
	option.UserID = userID
	err := template.MultiDelete(&option, reqMsg.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func Create(c *gin.Context) {
	var reqMsg template.CreateMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := parameter.GetUser(c)
	userID, _ := user["id"].(string)

	var option parameter.Option
	option.UserID = userID

	_, err := template.Create(&option, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func Get(c *gin.Context) {
	identify := c.Param("id")

	user := parameter.GetUser(c)
	userID, _ := user["id"].(string)
	var option parameter.Option
	option.UserID = userID

	m := template.FilterGet(&option, []string{"id"}, []interface{}{identify})

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"date": m,
	})
}

func Delete(c *gin.Context) {
	identify := c.Param("id")

	user := parameter.GetUser(c)
	userID, _ := user["id"].(string)
	var option parameter.Option
	option.UserID = userID

	err := template.Delete(&option, identify)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func Edit(c *gin.Context) {
	identify := c.Param("id")

	var reqMsg template.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := parameter.GetUser(c)
	userID, _ := user["id"].(string)
	var option parameter.Option
	option.UserID = userID

	err := template.Edit(&option, identify, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
