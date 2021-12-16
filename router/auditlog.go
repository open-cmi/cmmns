package router

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/controller/auditlog"
)

// AuditLogAuthGroup audit log group router
func AuditLogAuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/auditlog")
	{
		g.GET("/", auditlog.List)
	}
}
