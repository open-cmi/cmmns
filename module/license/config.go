package license

import (
	"encoding/json"
	"path"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/pkg/eyas"
)

var gConf Config

type Config struct {
	Lic        string `json:"lic"`
	PublicFile string `json:"public_file"`
}

func GetLicensePath() string {
	if gConf.Lic != "" {
		confDir := eyas.GetConfDir()
		return path.Join(confDir, "xsnos.lic")
	}
	return gConf.Lic
}

func GetPublicPemPath() string {
	if gConf.PublicFile != "" {
		return gConf.PublicFile
	}
	return path.Join(eyas.GetConfDir(), "public.pem")
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
