package manhour

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/module/manhour"
)

func List(c *gin.Context) {
	var option api.Option
	api.ParseParams(c, &option)

	count, results, err := manhour.List(&option)
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
	var reqMsg manhour.MultiDeleteMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	var option api.Option
	option.UserID = userID
	err := manhour.MultiDelete(&option, reqMsg.ID)
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
	var reqMsg manhour.CreateMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	var option api.Option
	option.UserID = userID

	_, err := manhour.Create(&option, &reqMsg)
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

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	var option api.Option
	option.UserID = userID

	m := manhour.Get(&option, "id", identify)

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"date": m,
	})
}

func Delete(c *gin.Context) {
	identify := c.Param("id")

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	var option api.Option
	option.UserID = userID

	err := manhour.Delete(&option, identify)
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

	var reqMsg manhour.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	var option api.Option
	option.UserID = userID

	err := manhour.Edit(&option, identify, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
