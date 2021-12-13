package agent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	model "github.com/open-cmi/cmmns/model/agent"
	agentmsg "github.com/open-cmi/cmmns/msg/agent"
	msg "github.com/open-cmi/cmmns/msg/common"
	"github.com/open-cmi/cmmns/utils"
)

// List list agents
func List(c *gin.Context) {

	var param msg.RequestQuery
	utils.ParseParams(c, &param)
	count, list, err := model.List(&param)
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

// CreateAgent create agent
func CreateAgent(c *gin.Context) {
	var createmsg agentmsg.CreateMsg
	if err := c.ShouldBindJSON(&createmsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 如果地址是localhost, 则转换成127.0.0.1, 目前只能支持ipv4地址
	if createmsg.Address == "localhost" {
		createmsg.Address = "127.0.0.1"
	}

	// 校验，这里暂时忽略
	err := model.CreateAgent(&createmsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	return
}

// DelAgent del agent
func DelAgent(c *gin.Context) {
	id := c.Param("id")
	err := model.DelAgent(id)
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

	taskid := fmt.Sprintf("deploy-task-%d", time.Now().Unix())
	var agents []model.Model = []model.Model{}
	for _, id := range dmsg.NodeID {
		mdl, err := model.GetAgent(id)
		if err != nil {
			continue
		}
		agents = append(agents, mdl)
	}

	err := Deploy(taskid, agents)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": taskid,
	})
	return
}
