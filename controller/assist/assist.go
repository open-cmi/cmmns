package assist

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/open-cmi/cmmns/config"
	"github.com/open-cmi/cmmns/modules/rassist"

	"gopkg.in/ini.v1"
)

func Enable(c *gin.Context) {
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
	for _, es := range conf.ExposedService {
		section, _ := file.NewSection(es.Name)
		section.NewKey("type", es.Type)
		section.NewKey("local_ip", es.LocalIP)
		section.NewKey("local_port", strconv.Itoa(int(es.LocalPort)))
		section.NewKey("remote_port", strconv.Itoa(int(es.RemotePort)))
	}
	file.SaveTo(frpCfg)

	err := rassist.RunClient(frpCfg)
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

func Disable(c *gin.Context) {
	rassist.Close()
	c.JSON(http.StatusOK, gin.H{
		"ret": 0,
		"msg": "",
	})
}
