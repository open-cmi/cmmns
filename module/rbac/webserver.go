package rbac

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac/middleware"
	"github.com/open-cmi/gobase/essential/webserver"
	"github.com/open-cmi/gobase/pkg/eyas"
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
	OptionAuthInit(e, gRbacMenuConf.Strict)
	MustAuthInit(e)

	e.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"ret": -1, "msg": "api not found"})
			return
		}

		pubdir := filepath.Join(eyas.GetWorkingDir(), "public")

		// try to serve static file
		fpath := filepath.Join(pubdir, path)
		if info, err := os.Stat(fpath); err == nil && !info.IsDir() {
			c.File(fpath)
			return
		}

		var indexFile string = "index.html"
		index := filepath.Join(pubdir, indexFile)
		if _, err := os.Stat(index); err == nil {
			c.File(index)
			return
		}

		c.String(http.StatusNotFound, "page not found")
	})
}

func init() {
	webserver.SetRBAC(&RBACAuth{})
}
