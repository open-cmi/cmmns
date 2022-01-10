package ticker

import (
	"github.com/open-cmi/cmmns/model/system"
)

// Init init start up
func Init() {
	go system.StartMonitor()
}
