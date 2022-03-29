package agent

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/common/errcode"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/agent"
	"github.com/open-cmi/cmmns/module/scheduler"
)

// KeepAlive keep alive
func KeepAlive(c *gin.Context) {
	// 获取device id
	devID := c.Query("dev_id")
	if devID == "" {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "dev id is required"})
		return
	}
	var option api.Option
	option.UserID = ""

	mdl := agent.Get("dev_id", devID)
	if mdl == nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": errcode.ErrNotRegistered,
			"msg": "agent not registered",
		})
		return
	}

	// if status is not online, udpate
	// update every 4 minute
	now := time.Now().Unix()
	if mdl.State != agent.StateOnline || now-mdl.UpdatedTime > int64(4*time.Minute) {
		mdl.State = agent.StateOnline
		err := mdl.Save()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": err.Error()})
			return
		}
	}

	// 先检查executor是否存在，如果不存在，则查询model
	sched := scheduler.GetScheduler()
	if sched == nil {
		logger.Errorf("sched is nil\n")
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "scheduler is nil"})
		return
	}

	consumer := sched.GetConsumer(devID)
	if consumer == nil {
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
}

func ReportResult(c *gin.Context) {
	// 获取device id
	devID := c.Query("dev_id")
	if devID == "" {
		c.JSON(http.StatusOK, gin.H{"ret": 1, "msg": "dev id is required"})
		return
	}
	// 解析结果数据
	var resp def.JobResponse
	if err := c.ShouldBindJSON(&resp); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
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

	// 根据job内容，修改结果，入库，删除缓存
	consumer.JobDone(&resp)

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

// Register 只有当数据库中不存在dev_id时，才需要注册
// 1. 数据库中，不存在任何该设备的信息，此时可以新建一个
// 2. 数据库中，有用户添加的设备信息，但是信息不全，此时，需要根据用户的地址信息去判断
// 2.1 当用户添加设备时，如果是程序部署的，则部署时会获取ip地址，根据内网外地址确定唯一设备
// 2.2 当用户添加设备时，不允许选择手动部署，即使手动部署，也按照不存在该设备处理
func Register(c *gin.Context) {
	clientIP := c.ClientIP()

	var msgobj agent.RegisterMsg
	if err := c.ShouldBindJSON(&msgobj); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	// 这里需要验证token

	// 根据外网地址与内网地址唯一确定一个agent
	mdl := agent.FilterGet(&api.Option{},
		[]string{"address", "local_address"},
		[]interface{}{clientIP, msgobj.LocalAddress},
	)
	if mdl != nil {
		mdl.HostName = msgobj.HostName
		mdl.DevID = msgobj.DevID
		mdl.State = agent.StateOnline
		mdl.Save()
		c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
		return
	}

	mdl = agent.New()
	mdl.DevID = msgobj.DevID
	mdl.LocalAddress = msgobj.LocalAddress
	mdl.HostName = msgobj.HostName
	mdl.ConnType = "manual"
	mdl.Address = clientIP
	mdl.State = agent.StateOnline
	mdl.Group = "root"
	mdl.Save()

	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}
