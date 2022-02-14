package assist

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/open-cmi/cmmns/essential/api"
	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/module/assist/infra"
	"github.com/open-cmi/cmmns/module/assist/router"
	"gopkg.in/ini.v1"
)

type RemoteService struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	LocalIP    string `json:"local_ip"`
	LocalPort  uint16 `json:"local_port"`
	RemotePort uint16 `json:"remote_port"`
}

type Config struct {
	ServerAddr string          `json:"server_addr"`
	ServerPort uint16          `json:"server_port"`
	Token      string          `json:"token,omitempty"`
	Service    []RemoteService `json:"services"`
}

var moduleConfig Config

func Init() error {
	// 根据配置文件，生成临时ini文件，然后传入参数
	tmpdir := os.TempDir()
	cfgFilePath := filepath.Join(tmpdir, "./frpc.ini")

	file := ini.Empty()
	comsec, _ := file.NewSection("common")
	comsec.NewKey("server_addr", moduleConfig.ServerAddr)
	comsec.NewKey("server_port", strconv.Itoa(int(moduleConfig.ServerPort)))
	if moduleConfig.Token != "" {
		comsec.NewKey("token", moduleConfig.Token)
	}
	for _, rs := range moduleConfig.Service {
		section, _ := file.NewSection(rs.Name)
		section.NewKey("type", rs.Type)
		section.NewKey("local_ip", rs.LocalIP)
		section.NewKey("local_port", strconv.Itoa(int(rs.LocalPort)))
		section.NewKey("remote_port", strconv.Itoa(int(rs.RemotePort)))
	}

	file.SaveTo(cfgFilePath)

	infra.Init(cfgFilePath)

	return nil
}

func init() {
	config.RegisterConfig("assist", &moduleConfig)
	api.RegisterAuthAPI("assist", router.AuthGroup)
}
