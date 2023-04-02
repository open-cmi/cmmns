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
	webserver.RegisterAuthRouter("agent", "/api/common/v3/agent")
	webserver.RegisterAuthAPI("agent", "GET", "/", List)
	webserver.RegisterAuthAPI("agent", "POST", "/", Create)
	webserver.RegisterAuthAPI("agent", "DELETE", "/:id", Delete)
	webserver.RegisterAuthAPI("agent", "PUT", "/:id", Edit)
	webserver.RegisterAuthAPI("agent", "POST", "/deploy/", Deploy)

	webserver.RegisterAuthRouter("master", "/api/common/v3/master/")
	webserver.RegisterAuthAPI("master", "GET", "/setting", GetSetting)
	webserver.RegisterAuthAPI("master", "GET", "/auto-master-setting", AutoGetMaster)
	webserver.RegisterAuthAPI("master", "POST", "/setting", EditSetting)

	webserver.RegisterUnauthRouter("agent", "/api/common/v3/agent")
	webserver.RegisterUnauthAPI("agent", "GET", "/get-job", GetJob)
	webserver.RegisterUnauthAPI("agent", "POST", "/report-result", ReportResult)
	webserver.RegisterUnauthAPI("agent", "GET", "/keep-alive", KeepAlive)
	webserver.RegisterUnauthAPI("agent", "POST", "/register", Register)
}
