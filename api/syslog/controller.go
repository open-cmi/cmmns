package syslog

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/syslog"
	"github.com/open-cmi/gobase/pkg/ginparam"
)

// List queries syslog records with pagination and filters
func List(c *gin.Context) {
	param := ginparam.ParseListOption(c)

	var filter syslog.QueryFilter
	if proc := c.Query("process"); proc != "" {
		filter.Process = proc
	}
	if procs := c.Query("processes"); procs != "" {
		for _, p := range strings.Split(procs, ",") {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				filter.Processes = append(filter.Processes, trimmed)
			}
		}
	}
	if keyword := c.Query("keyword"); keyword != "" {
		filter.Keyword = keyword
	}
	if level := c.Query("level"); level != "" {
		filter.Level = level
	}

	timeStartStr := c.Query("time_start")
	timeEndStr := c.Query("time_end")
	if timeStartStr != "" && timeEndStr != "" {
		timeStart, _ := strconv.ParseInt(timeStartStr, 10, 64)
		timeEnd, _ := strconv.ParseInt(timeEndStr, 10, 64)
		filter.TimeStart = timeStart
		filter.TimeEnd = timeEnd
	}

	count, list, err := syslog.QueryList(&param, &filter)
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
		},
	})
}

func init() {
	rbac.MustAuthAPI("syslog", "GET", "/", List,
		rbac.PermCode("cmmns:syslog", rbac.PermView),
	)
}
