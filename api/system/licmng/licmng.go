package licmng

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/licmng"
	"github.com/open-cmi/cmmns/module/user"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/webserver"
	"github.com/open-cmi/gobase/pkg/goparam"
)

// List list license
func QueryLicenseList(c *gin.Context) {

	query := goparam.ParseParams(c)

	userparam := goparam.GetUser(c)
	username := userparam["username"].(string)
	role := userparam["role"].(string)

	var filter licmng.QueryFilter
	if role != "admin" {
		filter.User = username
	}

	customer := c.Query("customer")
	if customer != "" {
		filter.Customer = customer
	}

	count, lics, err := licmng.QueryLicenseList(query, &filter)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": i18n.Sprintf("query license list failed"),
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

	if !licmng.IsEnable() {
		c.String(http.StatusNotFound, "404 Not Found")
		ah.InsertOperationLog(i18n.Sprintf("create license"), false)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("create license"), false)
		return
	}

	userparam := goparam.GetUser(c)
	username := userparam["username"].(string)
	role := userparam["role"].(string)

	// 如果
	if req.Version == "enterprise" && role != "admin" {
		ah.InsertOperationLog(i18n.Sprintf("create enterprise license"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "no permission"})
		return
	}

	_, err := licmng.CreateLicense(&req, username)
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
	param := goparam.ParseParams(c)

	usr := user.Get(param.UserID)
	if usr == nil {
		c.JSON(http.StatusForbidden, "")
		return
	}

	if param.Role != "admin" {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "no permission"})
		return
	}

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
	m := licmng.Get(id)
	if m == nil {
		ah.InsertOperationLog(i18n.Sprintf("download license"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "license not found"})
		return
	}

	content, err := licmng.CreateLicenseContent(m)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("download license"), false)

		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	fileName := fmt.Sprintf("%s.lic", m.MCode)

	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+fileName)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")

	_, err = c.Writer.Write([]byte(content))
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("download license"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -2, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("download license"), true)
}

func init() {
	webserver.RegisterMustAuthAPI("system", "GET", "/licmng/", QueryLicenseList)
	webserver.RegisterMustAuthAPI("system", "POST", "/licmng/", CreateLicense)
	webserver.RegisterMustAuthAPI("system", "GET", "/licmng/download/", DownloadLicense)
	webserver.RegisterMustAuthAPI("system", "DELETE", "/licmng/:id", DeleteLicense)
}
