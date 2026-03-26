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
	IgnoreLic bool              `json:"ignore_lic"`
	Roles     map[string][]Menu `json:"roles"`
}

var gRbacMenuConf RbacMenu

func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gRbacMenuConf)

	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(gRbacMenuConf)
	return v
}

func init() {
	config.RegisterConfig("rbac", Parse, Save)
}
