package setting

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/license"
	"github.com/open-cmi/cmmns/pkg/path"
	"github.com/open-cmi/cmmns/service/webserver"
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
	file, _ := c.FormFile("file")
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	defer src.Close()
	rootdir := path.GetRootPath()
	dst := filepath.Join(rootdir, "etc", "xsnos.lic")
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	os.Remove(file.Filename)
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func init() {
	webserver.RegisterAuthAPI("system", "GET", "/license/", GetLicenseInfo)
	webserver.RegisterAuthAPI("system", "POST", "/license/", UploadLicenseFile)
}
