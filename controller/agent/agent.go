package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/db"
	model "github.com/open-cmi/cmmns/model/agent"
	msg "github.com/open-cmi/cmmns/msg"
	agentmsg "github.com/open-cmi/cmmns/msg/agent"
	"github.com/open-cmi/cmmns/utils"
)

// List list agents
func List(c *gin.Context) {

	var param msg.RequestParams
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
	type PubMsg struct {
		TaskID string        `json:"taskid"`
		Data   []model.Model `json:"data"`
	}

	var pubmsg PubMsg = PubMsg{
		TaskID: taskid,
		Data:   agents,
	}
	msg, err := json.Marshal(pubmsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "convert to json failed"})
		return
	}

	db.GetCache(db.TaskCache).Publish(context.TODO(), "DeployAgent", msg)
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": taskid,
	})
	return
}
