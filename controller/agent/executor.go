package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/errcode"
	model "github.com/open-cmi/cmmns/model/agent"
	msg "github.com/open-cmi/cmmns/msg/agent"
	"github.com/open-cmi/cmmns/scheduler"
)

// KeepAlive keep alive
func KeepAlive(c *gin.Context) {
	// 获取device id
	devID := c.Query("dev_id")
	if devID == "" {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "dev id is required"})
		return
	}

	// 先检查executor是否存在，如果不存在，则查询model
	sched := scheduler.GetScheduler()
	if sched == nil {
		return
	}
	consumer := sched.GetConsumer(devID)
	if consumer == nil {
		var option model.Option
		option.Option.UserID = ""

		mdl := model.Get(&option, "dev_id", devID)
		if mdl == nil {
			c.JSON(http.StatusOK, gin.H{
				"ret": errcode.ErrNotRegistered,
				"msg": "agent not registered",
			})
			return
		}

		// 节点存在，需要更新信息
		if mdl.State != model.StateOnline {
			mdl.State = model.StateOnline
			err := mdl.Save()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": err.Error()})
				return
			}
		}

		consumer = sched.NewConsumer(&scheduler.ConsumerOption{
			Identity: devID,
			Group:    "default",
		})
		if consumer == nil {
			c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "register consumer failed"})
			return
		}
	}

	// 查看是否有配置变更

	// 查看是否有自己agent的任务
	if consumer.HasJob() {
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

// GetJob get job
func GetJob(c *gin.Context) {
	// 获取device id
	devID := c.Query("dev_id")
	if devID == "" {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "dev id is required"})
		return
	}

	sched := scheduler.GetScheduler()
	if sched == nil {
		return
	}

	consumer := sched.GetConsumer(devID)
	if consumer == nil {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "agent not exist"})
		return
	}

	job := consumer.GetJob()
	if job == nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": *job,
	})
	return
}

func Register(c *gin.Context) {
	clientIP := c.ClientIP()

	var msgobj msg.RegisterMsg
	if err := c.ShouldBindJSON(&msgobj); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	// 这里需要验证token

	// dd
	mdl := model.Get(&model.Option{}, "dev_id", msgobj.DevID)
	if mdl != nil {
		mdl.Address = clientIP
		mdl.HostName = msgobj.HostName
		mdl.LocalAddress = msgobj.LocalAddress
		mdl.Save()
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
		return
	}

	mdl = model.New()
	mdl.DevID = msgobj.DevID
	mdl.LocalAddress = msgobj.LocalAddress
	mdl.HostName = msgobj.HostName

	mdl.ConnType = "manual"
	mdl.Address = clientIP
	mdl.State = model.StateOnline
	mdl.Save()

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
	return
}
