package agentgroup

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/module/agent"
	"github.com/open-cmi/cmmns/module/agentgroup"
)

func List(c *gin.Context) {
	var option api.Option
	api.ParseParams(c, &option)

	count, results, err := agentgroup.List(&option)
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
	var reqMsg agentgroup.MultiDeleteMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	var option api.Option
	option.UserID = userID
	err := agentgroup.MultiDelete(&option, reqMsg.ID)
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
	var reqMsg agentgroup.CreateMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	var option api.Option
	option.UserID = userID

	_, err := agentgroup.Create(&option, &reqMsg)
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

	m := agentgroup.FilterGet(&option, []string{"id"}, []interface{}{identify})

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

	grp := agentgroup.Get(&option, "id", identify)
	if grp != nil {
		agentModel := agent.FilterGet(nil, []string{"group_name"}, []interface{}{grp.Name})
		if agentModel != nil {
			c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": i18n.Sprintf("group has been referenced")})
			return
		}
	}

	grp.Remove()

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
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
	var option api.Option
	option.UserID = userID

	err := agentgroup.Edit(&option, identify, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
