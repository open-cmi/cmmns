package auditlog

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/user"
	"github.com/open-cmi/cmmns/pkg/goparam"
)

func List(c *gin.Context) {
	param := goparam.ParseParams(c)

	usr := user.Get(param.UserID)
	if usr == nil {
		c.JSON(http.StatusForbidden, "")
		return
	}

	var paramnum int = 1
	var whereClause string
	var whereArgs []interface{}

	addr := c.Query("ip")
	if addr != "" {
		if whereClause != "" {
			whereClause += " and "
		}
		whereClause += fmt.Sprintf(`ip like %s`, sqldb.LikePlaceHolder(paramnum))
		whereArgs = append(whereArgs, addr)
		paramnum += 1
	}

	timeStartStr := c.Query("time_start")
	timeEndStr := c.Query("time_end")
	if timeStartStr != "" && timeEndStr != "" {
		if whereClause != "" {
			whereClause += " and "
		}
		timeStart, _ := strconv.Atoi(timeStartStr)
		timeEnd, _ := strconv.Atoi(timeEndStr)
		if timeStart < timeEnd {
			whereClause += fmt.Sprintf(`timestamp > %d and time < %d`, paramnum, paramnum+1)
			whereArgs = append(whereArgs, timeStart, timeEnd)
			paramnum += 2
		}
	}

	if paramnum != 1 {
		param.WhereClause = " where " + whereClause
		param.WhereArgs = whereArgs
	}

	count, list, err := auditlog.List(param)
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
