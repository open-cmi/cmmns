package time

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/sqldb"
	"github.com/open-cmi/cmmns/essential/webserver"
	"github.com/open-cmi/cmmns/pkg/goparam"
	"github.com/open-cmi/nays/module/object/time"
)

func ObjectNameList(c *gin.Context) {
	param := goparam.ParseParams(c)

	name := c.Query("name")
	if name != "" {
		param.WhereArgs = append(param.WhereArgs, name)
		param.WhereClause += fmt.Sprintf(" where name like %s", sqldb.LikePlaceHolder(len(param.WhereArgs)))
	}

	names, err := time.QueryTimeObjectNames(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": names,
	})
}

func QueryAbsoluteTimeObjectList(c *gin.Context) {
	param := goparam.ParseParams(c)
	name := c.Query("name")
	if name != "" {
		param.WhereArgs = append(param.WhereArgs, name)
		param.WhereClause += fmt.Sprintf(" where name like %s", sqldb.LikePlaceHolder(len(param.WhereArgs)))
	}

	count, results, err := time.QueryAbsoluteTimeObjectList(param)
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

func QueryPeriodTimeObjectList(c *gin.Context) {
	param := goparam.ParseParams(c)
	name := c.Query("name")
	if name != "" {
		param.WhereArgs = append(param.WhereArgs, name)
		param.WhereClause += fmt.Sprintf(" where name like %s", sqldb.LikePlaceHolder(len(param.WhereArgs)))
	}

	count, results, err := time.QueryPeriodTimeObjectList(param)
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

func CreateAbsoluteTime(c *gin.Context) {
	var req time.AbsoluteTimeObject
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	err := time.CreateAbsoluteTimeObject(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func CreatePeriodTime(c *gin.Context) {
	var req time.PeriodTimeObject
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := time.CreatePeriodTimeObject(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func DeleteTimeObject(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	err := time.DeleteTimeObject(req.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	webserver.RegisterAuthAPI("object", "GET", "/time/object-names/", ObjectNameList)
	webserver.RegisterAuthAPI("object", "GET", "/time/absolute/", QueryAbsoluteTimeObjectList)
	webserver.RegisterAuthAPI("object", "POST", "/time/absolute/", CreateAbsoluteTime)
	webserver.RegisterAuthAPI("object", "POST", "/time/absolute/del/", DeleteTimeObject)

	webserver.RegisterAuthAPI("object", "GET", "/time/period/", QueryPeriodTimeObjectList)
	webserver.RegisterAuthAPI("object", "POST", "/time/period/", CreatePeriodTime)
	webserver.RegisterAuthAPI("object", "POST", "/time/period/del/", DeleteTimeObject)
}
