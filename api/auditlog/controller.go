package auditlog

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/sshlogin"
	"github.com/open-cmi/gobase/pkg/ginparam"
)

func List(c *gin.Context) {
	param := ginparam.ParseListOption(c)

	var filter auditlog.QueryFilter
	addr := c.Query("ip")
	if addr != "" {
		filter.IP = addr
	}

	username := c.Query("username")
	if username != "" {
		filter.Username = username
	}

	timeStartStr := c.Query("time_start")
	timeEndStr := c.Query("time_end")
	if timeStartStr != "" && timeEndStr != "" {
		timeStart, _ := strconv.Atoi(timeStartStr)
		timeEnd, _ := strconv.Atoi(timeEndStr)
		filter.TimeStart = int64(timeStart)
		filter.TimeEnd = int64(timeEnd)
	}

	count, list, err := auditlog.QueryList(&param, &filter)
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

func ListSSHLogin(c *gin.Context) {
	status := c.DefaultQuery("status", "all")
	keyword := strings.TrimSpace(c.Query("keyword"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	count, results, err := sshlogin.QuerySSHLoginEvents(status, keyword, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": results,
		},
	})
}

func init() {
	rbac.MustAuthAPI("auditlog", "GET", "/", List,
		rbac.PermCode("cmmns:auditlog", rbac.PermView),
	)
	rbac.MustAuthAPI("auditlog", "GET", "/ssh-login/", ListSSHLogin,
		rbac.PermCode("cmmns:auditlog", rbac.PermView),
	)
}
