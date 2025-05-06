package rbac

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

type RbacMenu struct {
	NoLic []Menu            `json:"nolic"`
	Roles map[string][]Menu `json:"roles"`
}

var gRbacMenus RbacMenu

func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gRbacMenus)

	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(gRbacMenus)
	return v
}

func init() {
	config.RegisterConfig("rbac", Parse, Save)
}
