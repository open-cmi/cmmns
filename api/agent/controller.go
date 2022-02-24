package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/module/agent"
	"github.com/open-cmi/cmmns/module/auditlog"
)

// List list agents
func List(c *gin.Context) {
	var option api.Option
	api.ParseParams(c, &option)

	count, list, err := agent.List(&option)
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
}

// Create create agent
func Create(c *gin.Context) {
	var createmsg agent.CreateMsg
	if err := c.ShouldBindJSON(&createmsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 如果地址是localhost, 则转换成127.0.0.1, 目前只能支持ipv4地址
	if createmsg.Address == "localhost" {
		createmsg.Address = "127.0.0.1"
	}
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	// 校验，这里暂时忽略
	var option api.Option
	option.UserID = userID
	_, err := agent.Create(&option, &createmsg)
	if err != nil {
		auditlog.InsertLog(c, auditlog.OperationType, i18n.Sprintf("create agent failed"))
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	// 写日志操作
	auditlog.InsertLog(c, auditlog.OperationType, i18n.Sprintf("create agent successfully"))

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// Delete del agent
func Delete(c *gin.Context) {
	id := c.Param("id")
	user := api.GetUser(c)
	userID, _ := user["id"].(string)

	var option api.Option
	option.UserID = userID
	err := agent.Delete(&option, id)
	if err != nil {
		// 写日志操作
		auditlog.InsertLog(c, auditlog.OperationType, i18n.Sprintf("delete agent failed"))

		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 写日志操作
	auditlog.InsertLog(c, auditlog.OperationType, i18n.Sprintf("delete agent successfully"))

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// Deploy deploy agent
func Deploy(c *gin.Context) {
	var dmsg agent.DeployMsg
	if err := c.ShouldBindJSON(&dmsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// get user
	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	var option api.Option
	option.UserID = userID

	for _, id := range dmsg.ID {
		mdl := agent.Get(&option, []string{"id"}, []interface{}{id})
		if mdl == nil {
			continue
		}
		err := DeployRemote(mdl)
		if err != nil {
			mdl.State = agent.StateDeployFailed
		} else {
			mdl.State = agent.StateDeploySuccess
		}
		mdl.Save()
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func Edit(c *gin.Context) {
	identify := c.Param("id")

	var reqMsg agent.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	var option api.Option
	option.UserID = userID

	err := agent.Edit(&option, identify, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
