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

const maxBackgroundSize = 5 * 1024 * 1024 // 5MB

// GetBackground handles fetching the product background image from the server
func GetBackground(c *gin.Context) {
	bgPath, exists := prod.GetBackgroundPath()
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"ret": -1, "msg": "background file not found"})
		return
	}

	c.Header("Cache-Control", "public, max-age=86400")
	c.File(bgPath)
}

// UploadBackground handles uploading a custom background image
func UploadBackground(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	file, err := c.FormFile("file")
	if err != nil {
		file, err = c.FormFile("background")
	}
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload background"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "file is required"})
		return
	}

	if file.Size > maxBackgroundSize {
		ah.InsertOperationLog(i18n.Sprintf("upload background"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "file size exceeds 5MB limit"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowed[ext] {
		ah.InsertOperationLog(i18n.Sprintf("upload background"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": fmt.Sprintf("unsupported file extension %s", ext)})
		return
	}

	dst := filepath.Join(prod.GetBrandingDir(), "background"+ext)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload background"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	if err := prod.SetBackgroundExt(ext); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload background"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("upload background"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// DeleteBackground handles resetting the custom background
func DeleteBackground(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	if err := prod.DeleteBackground(); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("reset background"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("reset background"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	rbac.UnauthAPI("system", "GET", "/prod/background/", GetBackground)
	rbac.OptionAuthAPI("system", "POST", "/prod/background/upload/", UploadBackground, rbac.PermCode("cmmns:prodinfo", rbac.PermWrite))
	rbac.OptionAuthAPI("system", "DELETE", "/prod/background/", DeleteBackground, rbac.PermCode("cmmns:prodinfo", rbac.PermWrite))
}
