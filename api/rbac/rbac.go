package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/common/goparam"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/service/webserver"
)

func RoleNameList(c *gin.Context) {
	count, results, err := rbac.RoleNameList()
	if err != nil {
		c.JSON(200, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": results,
		},
	})
}

func RoleList(c *gin.Context) {
	var param goparam.Option

	err := goparam.ParseParams(c, &param)
	if err != nil {
		c.JSON(200, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	count, results, err := rbac.RoleList(&param)
	if err != nil {
		c.JSON(200, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": results,
		},
	})
}

func ModuleList(c *gin.Context) {
	var param goparam.Option

	err := goparam.ParseParams(c, &param)
	if err != nil {
		c.JSON(200, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	count, results, err := rbac.ModuleList(&param)
	if err != nil {
		c.JSON(200, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": results,
		},
	})
}

func init() {
	webserver.RegisterAuthRouter("rbac", "/api/rbac/v1")
	webserver.RegisterAuthAPI("rbac", "GET", "/role/", RoleList)
	webserver.RegisterAuthAPI("rbac", "GET", "/role/name-list/", RoleNameList)
	webserver.RegisterAuthAPI("rbac", "GET", "/module/", ModuleList)
}
