package rbac

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/gobase/essential/i18n"
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
		return
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
	rbac.RegisterMustAuthRouter("rbac", "/api/rbac/v1")
	rbac.MustAuthAPI("rbac", "GET", "/role/", RoleList, rbac.GetInitRoles())
	rbac.MustAuthAPI("rbac", "POST", "/role/delete/", DeleteRole, rbac.GetInitRoles())
	rbac.MustAuthAPI("rbac", "GET", "/role/name-list/", GetAllRoleNames, rbac.GetInitRoles())
}
