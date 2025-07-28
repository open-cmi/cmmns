package auditlog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/user"
	"github.com/open-cmi/gobase/pkg/goparam"
)

func List(c *gin.Context) {
	param := goparam.ParseParams(c)

	usr := user.Get(param.UserID)
	if usr == nil {
		c.JSON(http.StatusForbidden, "")
		return
	}

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

	count, list, err := auditlog.QueryList(param, &filter)
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
