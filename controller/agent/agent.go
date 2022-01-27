package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/controller"
	model "github.com/open-cmi/cmmns/model/agent"
	msg "github.com/open-cmi/cmmns/msg/agent"
)

// List list agents
func List(c *gin.Context) {
	var option model.Option
	controller.ParseParams(c, &option.Option)

	count, list, err := model.List(&option)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": list,
		}})
	return
}

// Create create agent
func Create(c *gin.Context) {
	var createmsg msg.CreateMsg
	if err := c.ShouldBindJSON(&createmsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 如果地址是localhost, 则转换成127.0.0.1, 目前只能支持ipv4地址
	if createmsg.Address == "localhost" {
		createmsg.Address = "127.0.0.1"
	}
	user := controller.GetUser(c)
	userID, _ := user["id"].(string)

	// 校验，这里暂时忽略
	var option model.Option
	option.Option.UserID = userID
	_, err := model.Create(&option, &createmsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	return
}

// Delete del agent
func Delete(c *gin.Context) {
	id := c.Param("id")
	user := controller.GetUser(c)
	userID, _ := user["id"].(string)

	var option model.Option
	option.Option.UserID = userID
	err := model.Delete(&option, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	return
}

// Deploy deploy agent
func Deploy(c *gin.Context) {
	var dmsg msg.DeployMsg
	if err := c.ShouldBindJSON(&dmsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// get user
	user := controller.GetUser(c)
	userID, _ := user["id"].(string)
	var option model.Option
	option.Option.UserID = userID

	for _, id := range dmsg.ID {
		agent := model.Get(&option, "id", id)
		if agent == nil {
			continue
		}
		var err error
		err = DeployRemote(agent)
		if err != nil {
			// 部署失败，写任务日志信息
			agent.State = model.StateDeployFailed
		} else {
			agent.State = model.StateDeploySuccess
		}
		agent.Save()
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
	option.Option.UserID = userID

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
