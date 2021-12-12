package agentgroup

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// List list agents
func List(c *gin.Context) {
	/*
		var param msg.RequestQuery
		utils.ParseParams(c, &param)
		count, list, err := model.List(&param)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
			return
		}*/

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   0,
			"results": []string{},
		}})
	return
}
