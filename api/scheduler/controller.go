package scheduler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/parameter"
	"github.com/open-cmi/cmmns/module/scheduler"
)

func List(c *gin.Context) {
	var option parameter.Option
	parameter.ParseParams(c, &option)

	count, results, err := scheduler.List(&option)
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
	var reqMsg scheduler.MultiDeleteMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := parameter.GetUser(c)
	userID, _ := user["id"].(string)

	var option parameter.Option
	option.UserID = userID
	err := scheduler.MultiDelete(&option, reqMsg.ID)
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

func Get(c *gin.Context) {
	identify := c.Param("id")

	user := parameter.GetUser(c)
	userID, _ := user["id"].(string)
	var option parameter.Option
	option.UserID = userID

	m := scheduler.Get(&option, "id", identify)

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

	err := scheduler.Delete(&option, identify)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
