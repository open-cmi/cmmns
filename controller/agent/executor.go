package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/open-cmi/cmmns/model/agent"
	"github.com/open-cmi/cmmns/scheduler"
)

// KeepAlive keep alive
func KeepAlive(c *gin.Context) {
	clientIP := c.ClientIP()

	// 获取device id
	deviceid := c.Query("deviceid")
	if deviceid == "" {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "need deviceid is required"})
		return
	}

	// 先查缓存是否存在
	executor, err := scheduler.GetExecutor(deviceid)
	if err != nil {
		mdl, err := model.GetAgentByAddress(clientIP)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "agent not exist"})
			return
		}
		// 新节点，需要查询数据进行更新
		err = model.UpdateDeviceID(clientIP, deviceid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": err.Error()})
			return
		}

		err = scheduler.RegisterExecutor(mdl.Name, deviceid, clientIP, 0)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "register executor failed"})
			return
		}
	} else {
		executor.Refresh()
	}

	// 查看是否有配置变更

	// 查看是否有自己agent的任务
	if scheduler.HasJob(&executor) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "",
			"data": map[string]string{
				"msgtype": "GetJob",
			}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	return
}

// GetSelfTask get self task
func GetSelfTask(c *gin.Context) {
	// 获取device id
	deviceid := c.Query("deviceid")
	if deviceid == "" {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "need deviceid is required"})
		return
	}

	executor, err := scheduler.GetExecutor(deviceid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "agent not exist"})
		return
	}

	task, err := scheduler.GetJob(&executor)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": "get task failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": task,
	})
	return
}
