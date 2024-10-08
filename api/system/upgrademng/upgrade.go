package upgrademng

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/module/system/upgrademng"
	"github.com/open-cmi/cmmns/pkg/eyas"
	"github.com/open-cmi/cmmns/service/webserver"
)

func UploadMetaFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	defer src.Close()
	content, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	os.Remove(file.Filename)

	// 校验一下内容
	var umi upgrademng.UpgradeMetaInfo
	err = json.Unmarshal(content, &umi)
	if err != nil {
		logger.Errorf("unmarshall failed: %s\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": "meta info file format invalid"})
		return
	}

	// 保存
	metaFile := fmt.Sprintf("%s.meta.json", umi.Prod)
	upgradeDir := filepath.Join(eyas.GetDataDir(), "upgrades")
	os.MkdirAll(upgradeDir, 0644)

	dst := filepath.Join(upgradeDir, metaFile)
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	defer out.Close()

	_, err = out.Write(content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": umi.Prod,
	})
}

func UploadPackage(c *gin.Context) {
	file, _ := c.FormFile("file")
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}
	defer src.Close()

	// 打开meta文件，根据meta文件来保存package
	// 保存
	upgradeDir := filepath.Join(eyas.GetDataDir(), "upgrades")
	os.MkdirAll(upgradeDir, 0644)

	dst := filepath.Join(upgradeDir, file.Filename)
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
	webserver.RegisterAuthAPI("system", "POST", "/upgrade-mng/upload-upgrade-meta/", UploadMetaFile)
	webserver.RegisterAuthAPI("system", "POST", "/upgrade-mng/upload-upgrade-package/", UploadPackage)
	webserver.RegisterAuthAPI("system", "POST", "/upgrade-mng/upgrade-start/", StartUpgrade)
}
