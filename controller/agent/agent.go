package agent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/controller"
	model "github.com/open-cmi/cmmns/model/agent"
	agentmsg "github.com/open-cmi/cmmns/msg/agent"
	msg "github.com/open-cmi/cmmns/msg/request"
	"github.com/open-cmi/cmmns/utils"
)

// List list agents
func List(c *gin.Context) {
	var param msg.RequestQuery
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
	var createmsg agentmsg.CreateMsg
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

// DeployAgent deploy agent
func DeployAgent(c *gin.Context) {
	var dmsg agentmsg.DeployMsg
	if err := c.ShouldBindJSON(&dmsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	user := controller.GetUser(c)
	userID, _ := user["id"].(string)

	taskid := fmt.Sprintf("deploy-task-%d", time.Now().Unix())
	var agents []model.Model = []model.Model{}
	for _, id := range dmsg.ID {
		mdl := model.Get(&model.ModelOption{
			UserID: userID,
		}, "id", id)
		if mdl == nil {
			continue
		}
		agents = append(agents, *mdl)
	}

	if len(agents) != 0 {
		err := Deploy(taskid, agents)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"ret": 1,
				"msg": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": taskid,
	})
	return
}
