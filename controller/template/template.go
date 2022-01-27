package template

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/controller"
	model "github.com/open-cmi/cmmns/model/template"
	msg "github.com/open-cmi/cmmns/msg/template"
)

func List(c *gin.Context) {
	var option model.Option
	controller.ParseParams(c, &option.Option)

	count, results, err := model.List(&option)
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
	return
}

func MultiDelete(c *gin.Context) {
	var reqMsg msg.MultiDeleteMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := controller.GetUser(c)
	userID, _ := user["id"].(string)

	var option model.Option
	option.UserID = userID
	err := model.MultiDelete(&option, reqMsg.Name)
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
	return
}

func Create(c *gin.Context) {
	var reqMsg msg.CreateMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := controller.GetUser(c)
	userID, _ := user["id"].(string)

	var option model.Option
	option.UserID = userID

	_, err := model.Create(&option, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
	return
}

func Get(c *gin.Context) {
	identify := c.Param("id")

	user := controller.GetUser(c)
	userID, _ := user["id"].(string)
	var option model.Option
	option.UserID = userID

	m := model.Get(&option, "id", identify)

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"date": m,
	})
	return
}

func Delete(c *gin.Context) {
	identify := c.Param("id")

	user := controller.GetUser(c)
	userID, _ := user["id"].(string)
	var option model.Option
	option.UserID = userID

	err := model.Delete(&option, identify)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
	return
}

func Edit(c *gin.Context) {
	identify := c.Param("id")

	var reqMsg msg.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := controller.GetUser(c)
	userID, _ := user["id"].(string)
	var option model.Option
	option.UserID = userID

	err := model.Edit(&option, identify, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
	return
}
