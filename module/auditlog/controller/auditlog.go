package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/module/auditlog/model"
)

func List(c *gin.Context) {
	var param api.Option
	api.ParseParams(c, &param)
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
