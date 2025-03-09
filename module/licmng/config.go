package licmng

import (
	"encoding/json"
	"path"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/pkg/eyas"
)

var gConf Config

type Config struct {
	PrivateFile string `json:"private_file"`
}

func GetPrivatePemPath() string {
	if gConf.PrivateFile != "" {
		return gConf.PrivateFile
	}
	return path.Join(eyas.GetConfDir(), "private.pem")
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
	config.RegisterConfig("licmng", Parse, Save)
}
