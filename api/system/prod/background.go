package prod

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/gobase/pkg/eyas"
)

// GetBackground handles fetching the product logo from the server
func GetBackground(c *gin.Context) {
	// Use eyas to get the standard data directory for the application
	logoPath := filepath.Join(eyas.GetDataDir(), "background.png")

	// Verify if the logo file actually exists on the filesystem
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		// Return a 404 if the logo hasn't been uploaded or set yet
		c.JSON(http.StatusNotFound, gin.H{"ret": -1, "msg": "background file not found"})
		return
	}

	// Serve the file directly with the appropriate content-type headers
	c.File(logoPath)
}

func init() {
	// Register the API as an unauthenticated route so it can be accessed by the login page
	// The path will be prefixed by the system's base path (e.g., /api/system/v1/prod/logo/)
	rbac.UnauthAPI("system", "GET", "/prod/background/", GetBackground)
}
