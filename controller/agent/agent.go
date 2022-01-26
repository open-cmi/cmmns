package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/controller"
	model "github.com/open-cmi/cmmns/model/agent"
	msg "github.com/open-cmi/cmmns/msg/agent"
	"github.com/open-cmi/cmmns/msg/request"
	"github.com/open-cmi/cmmns/utils"
)

// List list agents
func List(c *gin.Context) {
	var param request.RequestQuery
	utils.ParseParams(c, &param)
	user := controller.GetUser(c)
	userID, _ := user["id"].(string)

	count, list, err := model.List(&model.ModelOption{
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
	_, err := model.Create(&model.ModelOption{
		UserID: userID,
	}, &createmsg)
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

	err := model.Delete(&model.ModelOption{
		UserID: userID,
	}, id)
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

	for _, id := range dmsg.ID {
		agent := model.Get(&model.ModelOption{
			UserID: userID,
		}, "id", id)
		if agent == nil {
			continue
		}
		var err error
		if agent.IsLocal {
			err = DeployLocal()
		} else {
			err = DeployRemote(agent)
		}
		if err != nil {
			// 部署失败，写任务日志信息
			agent.Reason = err.Error()
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
