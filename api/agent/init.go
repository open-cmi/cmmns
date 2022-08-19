package agent

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/service/webserver"
)

var gConf Config

func Init(confmsg json.RawMessage) error {
	err := json.Unmarshal(confmsg, &gConf)
	return err
}

func Save() json.RawMessage {
	vb, _ := json.Marshal(&gConf)
	return vb
}

func init() {
	config.RegisterConfig("cluster", Init, Save)
	webserver.RegisterAuthRouter("agent", AuthGroup)
	webserver.RegisterUnauthRouter("agent", UnauthGroup)
}
