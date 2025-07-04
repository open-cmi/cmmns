package license

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

var gConf Config

type Config struct {
	Lic        string `json:"lic"`
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
