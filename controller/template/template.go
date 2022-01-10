package template

import (
	"net/http"

	"github.com/gin-gonic/gin"

	model "github.com/open-cmi/cmmns/model/template"
	"github.com/open-cmi/cmmns/msg/request"
	msg "github.com/open-cmi/cmmns/msg/template"
	"github.com/open-cmi/cmmns/utils"
)

func List(c *gin.Context) {
	var param request.RequestQuery
	utils.ParseParams(c, &param)
	count, results, err := model.List(&param)
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
	_, err := model.Create(&reqMsg)
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

	m := model.Get(identify)
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"date": m,
	})
	return
}

func Delete(c *gin.Context) {
	identify := c.Param("id")

	err := model.Delete(identify)
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
	err := model.Edit(identify, &reqMsg)
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
