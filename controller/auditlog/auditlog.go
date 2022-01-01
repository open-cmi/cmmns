package auditlog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/open-cmi/cmmns/model/auditlog"
	commonmsg "github.com/open-cmi/cmmns/msg/request"
	"github.com/open-cmi/cmmns/utils"
)

func List(c *gin.Context) {
	var param commonmsg.RequestQuery
	utils.ParseParams(c, &param)
	count, list, err := model.List(&param)
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
