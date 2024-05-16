package rbac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/pkg/goparam"
	"github.com/open-cmi/cmmns/service/webserver"
)

func RoleNameList(c *gin.Context) {
	count, results, err := rbac.RoleNameList()
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

func RoleList(c *gin.Context) {
	param := goparam.ParseParams(c)

	count, results, err := rbac.RoleList(param)
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

func DeleteRole(c *gin.Context) {
	ID := c.Param("id")

	user := goparam.GetUser(c)
	userID, _ := user["id"].(string)
	err := rbac.DeleteRole(&goparam.Param{
		UserID: userID,
	}, ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func init() {
	webserver.RegisterAuthRouter("rbac", "/api/rbac/v1")
	webserver.RegisterAuthAPI("rbac", "GET", "/role/", RoleList)
	webserver.RegisterAuthAPI("rbac", "DELETE", "/role/:id", DeleteRole)
	webserver.RegisterAuthAPI("rbac", "GET", "/role/name-list/", RoleNameList)
}
