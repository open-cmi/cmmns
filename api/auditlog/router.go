package auditlog

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/service/webserver"
)

// AuthGroup audit log group router
func AuthGroup(e *gin.Engine) {

	g := e.Group("/api/common/v3/auditlog")
	{
		g.GET("/", List)
	}
}

func init() {
	webserver.RegisterAuthRouter("auditlog", AuthGroup)
}
