package secretkey

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/module/secretkey"
)

func List(c *gin.Context) {
	var param api.Option
	api.ParseParams(c, &param)

	count, results, err := secretkey.List(&param)
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
	var reqMsg secretkey.MultiDeleteMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	var option api.Option
	option.UserID = userID
	err := secretkey.MultiDelete(&option, reqMsg.Name)
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
	var reqMsg secretkey.CreateMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	_, err := secretkey.Create(&api.Option{
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

	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	m := secretkey.Get(&api.Option{
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

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	err := secretkey.Delete(&api.Option{
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

	var reqMsg secretkey.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	err := secretkey.Edit(&api.Option{
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
