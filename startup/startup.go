package startup

import (
	"github.com/open-cmi/cmmns/model/systeminfo"
	"github.com/open-cmi/cmmns/plugins"
)

// Init init start up
func Init() {
	systeminfo.InitDeviceID()
	go systeminfo.StartMonitor()
	go plugins.Init()
}
