package time

import "github.com/open-cmi/cmmns/initial"

func Init() error {
	go ObjectEventLoop()
	return nil
}

func init() {
	initial.Register("time-object", initial.DefaultPriority, Init)
}
