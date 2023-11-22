package auditlog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/user"
)

func List(c *gin.Context) {
	param := goparam.ParseParams(c)

	usr := user.Get("id", param.UserID)
	if usr == nil || !usr.HasReadPermision("Audit log") {
		c.JSON(http.StatusForbidden, "")
		return
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
