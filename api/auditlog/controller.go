package auditlog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/common/parameter"
	"github.com/open-cmi/cmmns/module/auditlog"
)

func List(c *gin.Context) {
	var param parameter.Option
	parameter.ParseParams(c, &param)
	count, list, err := auditlog.List(&param)
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
	return
}
