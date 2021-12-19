package ticker

import (
	"github.com/open-cmi/cmmns/model/systeminfo"
)

// Init init start up
func Init() {
	go systeminfo.StartMonitor()
}
