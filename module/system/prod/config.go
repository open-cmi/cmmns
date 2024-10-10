package prod

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

type Menu struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	Icon         string `json:"icon,omitempty"`
	Require      bool   `json:"require,omitempty"`
	Experimental bool   `json:"experimental,omitempty"`
	Children     []Menu `json:"children,omitempty"`
}

type NavConfig struct {
	Experimental bool   `json:"experimental"`
	Menus        []Menu `json:"menus"`
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
