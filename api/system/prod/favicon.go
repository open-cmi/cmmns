package prod

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/system/prod"
	"github.com/open-cmi/gobase/essential/i18n"
)

const maxFaviconSize = 512 * 1024 // 512KB

// GetFavicon handles fetching the browser favicon from the server
func GetFavicon(c *gin.Context) {
	favPath, exists := prod.GetFaviconPath()
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"ret": -1, "msg": "favicon file not found"})
		return
	}

	c.Header("Cache-Control", "public, max-age=86400")
	c.File(favPath)
}

// UploadFavicon handles uploading a custom browser favicon
func UploadFavicon(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	file, err := c.FormFile("file")
	if err != nil {
		file, err = c.FormFile("favicon")
	}
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload favicon"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "file is required"})
		return
	}

	if file.Size > maxFaviconSize {
		ah.InsertOperationLog(i18n.Sprintf("upload favicon"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "file size exceeds 512KB limit"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".ico": true,
		".png": true,
		".svg": true,
	}
	if !allowed[ext] {
		ah.InsertOperationLog(i18n.Sprintf("upload favicon"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": fmt.Sprintf("unsupported file extension %s", ext)})
		return
	}

	dst := filepath.Join(prod.GetBrandingDir(), "favicon"+ext)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload favicon"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	if err := prod.SetFaviconExt(ext); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload favicon"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("upload favicon"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// DeleteFavicon handles resetting the custom favicon
func DeleteFavicon(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	if err := prod.DeleteFavicon(); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("reset favicon"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("reset favicon"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	rbac.UnauthAPI("system", "GET", "/prod/favicon/", GetFavicon)
	rbac.OptionAuthAPI("system", "POST", "/prod/favicon/upload/", UploadFavicon, rbac.PermCode("cmmns:prodinfo", rbac.PermWrite))
	rbac.OptionAuthAPI("system", "DELETE", "/prod/favicon/", DeleteFavicon, rbac.PermCode("cmmns:prodinfo", rbac.PermWrite))
}
