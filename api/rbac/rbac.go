package rbac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/webserver"
	"github.com/open-cmi/gobase/pkg/goparam"
)

func GetAllRoleNames(c *gin.Context) {
	names, err := rbac.GetAllRoleNames()
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
	var apimsg rbac.RoleDeleteRequest
	if err := c.ShouldBindJSON(&apimsg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	user := goparam.GetUser(c)
	curRole, _ := user["role"].(string)
	if curRole != "admin" {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": i18n.Sprintf("no permission")})
	}

	err := rbac.DeleteRole(apimsg.ID)
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
	webserver.RegisterMustAuthAPI("rbac", "POST", "/role/delete/", DeleteRole)
	webserver.RegisterMustAuthAPI("rbac", "GET", "/role/name-list/", GetAllRoleNames)
}
