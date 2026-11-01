package setting

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
)

// GetFromRequest dynamically parses schema and host (including port if present) from http.Request
func GetFromRequest(req *http.Request) string {
	schema := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		schema = "https"
	}

	host := req.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = req.Host
	}

	return fmt.Sprintf("%s:%s", schema, host)
}

func GetPublicNet(c *gin.Context) {
	pn := GetFromRequest(c.Request)
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": pn,
	})
}

func init() {
	rbac.OptionAuthAPI("system", "GET", "/pubnet/", GetPublicNet)
}
