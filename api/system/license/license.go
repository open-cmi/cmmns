package setting

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/auditlog"
	"github.com/open-cmi/cmmns/module/license"
	"github.com/open-cmi/gobase/essential/events"
	"github.com/open-cmi/gobase/essential/i18n"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/webserver"
)

func GetLicenseInfo(c *gin.Context) {

	licInfo, err := license.GetLicenseInfo()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 返回LAN参数
	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": licInfo,
	})
}

func UploadLicenseFile(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)
	file, _ := c.FormFile("file")
	src, err := file.Open()
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload license file"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	defer src.Close()
	content, err := io.ReadAll(src)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload license file"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	os.Remove(file.Filename)
	err = license.VerifyLicenseContent(string(content))
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload license file"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	dst := license.GetLicensePath()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload license file"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	defer out.Close()

	_, err = out.Write(content)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("upload license file"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("upload license file"), true)
	events.Notify("check-license-valid", nil)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func SetProductSerial(c *gin.Context) {
	ah := auditlog.NewAuditHandler(c)

	var req struct {
		Serial  string `json:"serial"`
		Product string `json:"product"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		ah.InsertOperationLog(i18n.Sprintf("set product serial"), false)
		return
	}
	if req.Product == "" {
		req.Product = "swapi"
	}
	logger.Debugf("set product serial: %s, %s", req.Serial, req.Product)
	err := license.SetProductSerial(req.Serial, req.Product)
	if err != nil {
		ah.InsertOperationLog(i18n.Sprintf("set product serial"), false)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	ah.InsertOperationLog(i18n.Sprintf("set product serial"), true)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/license/", GetLicenseInfo)
	webserver.RegisterAuthAPI("system", "POST", "/license/upload/", UploadLicenseFile)
	webserver.RegisterUnauthAPI("system", "POST", "/license/product-serial/", SetProductSerial)
}
