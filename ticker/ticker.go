package ticker

import (
	"github.com/open-cmi/cmmns/model/system"

	"github.com/robfig/cron/v3"
)

// Init init start up
func Init() {
	c := cron.New(cron.WithSeconds())

	// run every minitue
	c.AddFunc("0 */5 * * * *", func() {
		system.StartMonitor()
	})
	c.Start()
}
