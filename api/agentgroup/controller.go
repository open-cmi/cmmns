package agentgroup

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/module/agentgroup"
)

func List(c *gin.Context) {
	var param api.Option
	api.ParseParams(c, &param)

	count, results, err := agentgroup.List(&param)
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
	var reqMsg agentgroup.MultiDeleteMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	var modelOption api.Option
	modelOption.UserID = userID
	err := agentgroup.MultiDelete(&modelOption, reqMsg.ID)
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
	var reqMsg agentgroup.CreateMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	_, err := agentgroup.Create(&api.Option{
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

	m := agentgroup.Get(&api.Option{
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
	err := agentgroup.Delete(&api.Option{
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

	var reqMsg agentgroup.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	err := agentgroup.Edit(&api.Option{
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
