package license

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
)

var gConf Config

type Config struct {
	PublicFile string `json:"public_file"`
}

func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gConf)
	return err
}

func Save() json.RawMessage {
	v, _ := json.Marshal(gConf)
	return v
}

func init() {
	config.RegisterConfig("license", Parse, Save)
}
