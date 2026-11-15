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

const maxLogoSize = 2 * 1024 * 1024 // 2MB

// GetLogo handles fetching the product logo from the server
func GetLogo(c *gin.Context) {
	logoPath, exists := prod.GetLogoPath()
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"ret": -1, "msg": "logo file not found"})
		return
	}

	c.Header("Cache-Control", "public, max-age=86400")
	c.File(logoPath)
}

// UploadLogo handles uploading a custom product logo
func UploadLogo(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	file, err := c.FormFile("file")
	if err != nil {
		file, err = c.FormFile("logo")
	}
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload logo"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "file is required"})
		return
	}

	if file.Size > maxLogoSize {
		ah.InsertOperationLog(i18n.Sprintf("upload logo"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "file size exceeds 2MB limit"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".png":  true,
		".svg":  true,
		".jpg":  true,
		".jpeg": true,
		".webp": true,
	}
	if !allowed[ext] {
		ah.InsertOperationLog(i18n.Sprintf("upload logo"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": fmt.Sprintf("unsupported file extension %s", ext)})
		return
	}

	dst := filepath.Join(prod.GetBrandingDir(), "logo"+ext)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload logo"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	if err := prod.SetLogoExt(ext); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload logo"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("upload logo"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

// DeleteLogo handles resetting the custom logo
func DeleteLogo(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	if err := prod.DeleteLogo(); err != nil {
		ah.InsertOperationLog(i18n.Sprintf("reset logo"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("reset logo"), true)
	c.JSON(http.StatusOK, gin.H{"ret": 0, "msg": ""})
}

func init() {
	rbac.UnauthAPI("system", "GET", "/prod/logo/", GetLogo)
	rbac.OptionAuthAPI("system", "POST", "/prod/logo/upload/", UploadLogo, rbac.PermCode("cmmns:prodinfo", rbac.PermWrite))
	rbac.OptionAuthAPI("system", "DELETE", "/prod/logo/", DeleteLogo, rbac.PermCode("cmmns:prodinfo", rbac.PermWrite))
}
