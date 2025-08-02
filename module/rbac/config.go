package rbac

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

type Menu struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Icon     string `json:"icon,omitempty"`
	Children []Menu `json:"children,omitempty"`
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
