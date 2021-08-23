package startup

import (
	"github.com/open-cmi/cmmns/model/systeminfo"
)

// Init init start up
func Init() {
	systeminfo.InitDeviceID()
	go systeminfo.StartMonitor()
}
