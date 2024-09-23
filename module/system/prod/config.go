package prod

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type Menu struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Icon     string `json:"icon,omitempty"`
	Hidden   bool   `json:"hidden,omitempty"`
	Require  bool   `json:"require,omitempty"`
	Children []Menu `json:"children,omitempty"`
}

type NavConfig struct {
	Menus []Menu `json:"menus"`
}

var gNavConf NavConfig

func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gNavConf)

	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(gNavConf)
	return v
}

func init() {
	config.RegisterConfig("nav", Parse, Save)
}
