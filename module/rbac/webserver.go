package rbac

import (
	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac/middleware"
	"github.com/open-cmi/gobase/essential/webserver"
)

type RBACAuth struct {
}

func (i *RBACAuth) Init(e *gin.Engine) {
	// init webserver
	middleware.DefaultMiddleware(e)

	// workDir := eyas.GetWorkingDir()
	// dir := fmt.Sprintf("%s/static/", workDir)
	// e.Static("/api-static/", dir)

	middleware.SessionMiddleware(e)
	middleware.JWTMiddleware(e)
	UnauthInit(e)
	if !gRbacMenuConf.Strict {
		OptionAuthInit(e)
	}
	middleware.AuthMiddleware(e)
	if gRbacMenuConf.Strict {
		OptionAuthInit(e)
	}
	MustAuthInit(e)
}

func init() {
	webserver.SetRBAC(&RBACAuth{})
}
