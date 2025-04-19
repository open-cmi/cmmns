package rbac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/webserver"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/pkg/goparam"
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
	webserver.RegisterMustAuthRouter("rbac", "/api/rbac/v1")
	webserver.RegisterMustAuthAPI("rbac", "GET", "/role/", RoleList)
	webserver.RegisterMustAuthAPI("rbac", "DELETE", "/role/:id", DeleteRole)
	webserver.RegisterMustAuthAPI("rbac", "GET", "/role/name-list/", RoleNameList)
}
