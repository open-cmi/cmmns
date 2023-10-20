package status

import (
	"net/http"

	"github.com/open-cmi/cmmns/common/parameter"
	"github.com/open-cmi/cmmns/module/system"

	"github.com/gin-gonic/gin"
)

// History list status history info
func StatusHistoryList(c *gin.Context) {
	var option parameter.Option
	parameter.ParseParams(c, &option)

	count, list, err := system.List(&option)
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

// Get get device info
func StatusHistoryGet(c *gin.Context) {
	id := c.Param("id")
	user := parameter.GetUser(c)
	userID, _ := user["id"].(string)

	var option parameter.Option
	option.UserID = userID
	mdl := system.Get(&option, "id", id)
	if mdl == nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "not exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": "", "data": mdl})
}
