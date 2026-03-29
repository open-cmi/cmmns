package powermanagement

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/system/powermanagement"
	"github.com/open-cmi/gobase/essential/i18n"
)

func ListSchedules(c *gin.Context) {
	schedules, err := powermanagement.GetSchedules()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": schedules})
}

func CreateSchedule(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	var task powermanagement.Schedule
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	err := powermanagement.AddSchedule(&task)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("create scheduled power task"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("create scheduled power task"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func DeleteSchedule(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	var body struct {
		ID string `json:"id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "id is required"})
		return
	}

	err := powermanagement.DeleteSchedule(body.ID)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("delete scheduled power task"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("delete scheduled power task"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	rbac.OptionAuthAPI("system", "GET", "/power-management/schedule/", ListSchedules, rbac.GetInitRoles())
	rbac.OptionAuthAPI("system", "POST", "/power-management/schedule/", CreateSchedule, rbac.GetInitRoles())
	rbac.OptionAuthAPI("system", "POST", "/power-management/schedule/delete/", DeleteSchedule, rbac.GetInitRoles())
}
