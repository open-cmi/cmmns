package router

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog/controller"
)

// AuthGroup audit log group router
func AuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/auditlog")
	{
		g.GET("/", controller.List)
	}
}
