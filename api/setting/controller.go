package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/open-cmi/cmmns/common/api"
	"github.com/open-cmi/cmmns/module/setting"
)

func List(c *gin.Context) {
	var option api.Option
	api.ParseParams(c, &option)

	count, results, err := setting.List(&option)
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

func Get(c *gin.Context) {
	identify := c.Param("id")

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	var option api.Option
	option.UserID = userID

	m := setting.FilterGet(&option, []string{"id"}, []interface{}{identify})

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"date": m,
	})
}

func Edit(c *gin.Context) {
	identify := c.Param("id")

	var reqMsg setting.EditMsg

	if err := c.ShouldBindJSON(&reqMsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := api.GetUser(c)
	userID, _ := user["id"].(string)
	var option api.Option
	option.UserID = userID

	err := setting.Edit(&option, identify, &reqMsg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
