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
	if curRole == apimsg.ID {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": i18n.Sprintf("no permission")})
		return
	}

	role := rbac.GetByID(apimsg.ID)
	if role == nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": i18n.Sprintf("role is not exist")})
		return
	}

	if role.Name == curRole {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": i18n.Sprintf("the role can not be deleted")})
		return
	}

	if role.Preset {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": i18n.Sprintf("preset role can not be deleted")})
		return
	}

	err := role.Remove()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func CreateRole(c *gin.Context) {
	var apimsg struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions"`
	}
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

	err := rbac.CreateRole(apimsg.Name, apimsg.Description, apimsg.Permissions)
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
	rbac.MustAuthAPI("rbac", "GET", "/role/", RoleList, rbac.PermCode("rbac", rbac.PermView))
	rbac.MustAuthAPI("rbac", "POST", "/role/", CreateRole, rbac.PermCode("rbac", rbac.PermWrite))
	rbac.MustAuthAPI("rbac", "POST", "/role/create/", CreateRole, rbac.PermCode("rbac", rbac.PermWrite))
	rbac.MustAuthAPI("rbac", "POST", "/role/delete/", DeleteRole, rbac.PermCode("rbac", rbac.PermDelete))
	rbac.MustAuthAPI("rbac", "GET", "/role/name-list/", GetAllRoleNames, rbac.PermCode("rbac", rbac.PermView))
}
