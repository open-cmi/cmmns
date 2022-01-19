package assist

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/config"
	massist "github.com/open-cmi/cmmns/modules/assist"

	"gopkg.in/ini.v1"
)

func GetAssist(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"ret":  0,
		"msg":  "",
		"data": massist.IsRunning(),
	})
}

func SetAssist(c *gin.Context) {

	var msg struct {
		Enable bool `json:"enable"`
	}

	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret": -1, "msg": err.Error()})
		return
	}

	if !msg.Enable {
		massist.Close()
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "",
		})
		return
	}

	// 根据配置文件，生成临时ini文件，然后传入参数
	tmpdir := os.TempDir()
	frpCfg := filepath.Join(tmpdir, "./frpc.ini")

	file := ini.Empty()
	comsec, _ := file.NewSection("common")
	conf := config.GetConfig().RemoteAssist
	comsec.NewKey("server_addr", conf.ServerAddr)
	comsec.NewKey("server_port", strconv.Itoa(int(conf.ServerPort)))
	if conf.Token != "" {
		comsec.NewKey("token", conf.Token)
	}
	for _, rs := range conf.Service {
		section, _ := file.NewSection(rs.Name)
		section.NewKey("type", rs.Type)
		section.NewKey("local_ip", rs.LocalIP)
		section.NewKey("local_port", strconv.Itoa(int(rs.LocalPort)))
		section.NewKey("remote_port", strconv.Itoa(int(rs.RemotePort)))
	}
	file.SaveTo(frpCfg)

	err := massist.RunClient(frpCfg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 1,
			"msg": "start assist service failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
