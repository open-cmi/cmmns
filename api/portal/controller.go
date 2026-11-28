package portal

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/portal"
	"github.com/open-cmi/cmmns/module/rbac"
)

// List returns portal product entries scanned from /opt/*/etc/eyas-portal.json
func List(c *gin.Context) {
	products, err := portal.ScanProducts("")
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"ret": -1, "msg": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": products,
	})
}

func init() {
	rbac.MustAuthAPI("portal", "GET", "/", List)
}
