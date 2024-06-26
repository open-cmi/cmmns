package prod

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type SubMenu struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type Menu struct {
	Path     string    `json:"path"`
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	Children []SubMenu `json:"children"`
}

type ProdBriefInfo struct {
	Name   string `json:"name"`
	Footer string `json:"footer"`
}

type ProdConfig struct {
	Name   string `json:"name"`
	Footer string `json:"footer"`
	Nav    []Menu `json:"nav"`
}

var gProdConf ProdConfig

func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gProdConf)
	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(gProdConf)
	return v
}

func init() {
	config.RegisterConfig("prod", Parse, Save)
}
