package upgrademng

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/module/rbac"
	"github.com/open-cmi/cmmns/module/system/upgrademng"
	"github.com/open-cmi/gobase/pkg/eyas"
)

func UploadPackage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	defer src.Close()

	// 打开meta文件，根据meta文件来保存package
	// 保存
	upgradeDir := filepath.Join(eyas.GetDataDir(), "upgrades")
	os.MkdirAll(upgradeDir, 0755)

	dst := filepath.Join(upgradeDir, filepath.Base(file.Filename))
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	written, err := io.Copy(out, src)
	_ = out.Close()
	if err != nil {
		_ = os.Remove(dst)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	if file.Size > 0 && written != file.Size {
		_ = os.Remove(dst)
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "upload incomplete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": file.Filename,
	})
}

func StartUpgrade(c *gin.Context) {
	var req upgrademng.UpgradeRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	// 发送升级指令
	err := upgrademng.StartUpgrade(&req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}

func init() {
	rbac.OptionAuthAPI("system", "POST", "/upgrade-mng/upload-upgrade-package/", UploadPackage,
		rbac.PermCode("cmmns:upgrade", rbac.PermWrite),
	)
	rbac.OptionAuthAPI("system", "POST", "/upgrade-mng/upgrade-start/", StartUpgrade,
		rbac.PermCode("cmmns:upgrade", rbac.PermWrite),
	)
}
