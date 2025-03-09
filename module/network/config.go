package network

import (
	"encoding/json"
	"runtime"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/service/initial"
)

var gConf Config

type Config struct {
	Engine   string   `json:"engine,omitempty"`
	ConfFile string   `json:"conf_file"`
	Devices  []string `json:"devices"`
}

func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	if gConf.Engine == "" {
		gConf.Engine = "netplan"
	}

	return err
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func Init() error {
	// 如果配置文件不包含管理口，则从数据库中获取
	if len(gConf.Devices) == 0 {
		gConf.Devices = LoadNetworkManagementInterface()
	}
	if len(gConf.Devices) != 0 {
		return NetworkApply(&gConf)
	}
	return nil
}

func init() {
	config.RegisterConfig("network", Parse, Save)
	if runtime.GOOS == "linux" {
		initial.Register("network", initial.DefaultPriority, Init)
	}
}
