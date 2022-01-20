package template

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/controller"
	model "github.com/open-cmi/cmmns/model/template"
	"github.com/open-cmi/cmmns/msg/request"
	msg "github.com/open-cmi/cmmns/msg/template"
	"github.com/open-cmi/cmmns/utils"
)

func List(c *gin.Context) {
	var param request.RequestQuery
	utils.ParseParams(c, &param)

	user := controller.GetUser(c)
	userID, _ := user["id"].(string)

	count, results, err := model.List(&model.ModelOption{
		UserID: userID,
	}, &param)
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

	_, err := model.Create(&model.ModelOption{
		UserID: userID,
	}, &reqMsg)
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

	m := model.Get(&model.ModelOption{
		UserID: userID,
	}, identify)

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
	err := model.Delete(&model.ModelOption{
		UserID: userID,
	}, identify)
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
	err := model.Edit(&model.ModelOption{
		UserID: userID,
	}, identify, &reqMsg)
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
