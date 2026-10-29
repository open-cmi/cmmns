package rbac

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

type RbacConf struct {
	Strict    bool `json:"strict"`
	IgnoreLic bool `json:"ignore_lic"`
}

var gRbacMenuConf RbacConf

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
