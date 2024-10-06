package licmng

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/i18n"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/licmng"
	"github.com/open-cmi/cmmns/pkg/goparam"
	"github.com/open-cmi/cmmns/service/webserver"
)

// List list license
func ListLicense(c *gin.Context) {

	query := goparam.ParseParams(c)

	count, lics, err := licmng.ListLicense(query)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("list license failed"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
		"data": map[string]interface{}{
			"count":   count,
			"results": lics,
		},
	})
}

func CreateLicense(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	var req licmng.CreateLicenseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("create license"), false)
		return
	}
	err := licmng.CreateLicense(&req)
	if err != nil {

		ah.InsertOperationLog(i18n.Sprintf("create license"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("create license"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func DeleteLicense(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "license should not be empty"})
		ah.InsertOperationLog(i18n.Sprintf("delete license"), false)
		return
	}

	err := licmng.DeleteLicense(id)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("delete license"), false)

		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	ah.InsertOperationLog(i18n.Sprintf("delete license"), true)

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func DownloadLicense(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	id := c.Query("id")
	content, err := licmng.CreateLicenseContent(id)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("download license"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	fileName := fmt.Sprintf("%s.lic", id)

	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+fileName)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	_, err = c.Writer.Write([]byte(content))
	if err != nil {

		ah.InsertOperationLog(i18n.Sprintf("download license"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("download license"), true)
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/licmng/", ListLicense)
	webserver.RegisterAuthAPI("system", "POST", "/licmng/", CreateLicense)
	webserver.RegisterAuthAPI("system", "GET", "/licmng/download/", DownloadLicense)
	webserver.RegisterAuthAPI("system", "DELETE", "/licmng/:id", DeleteLicense)
}
