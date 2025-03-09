package events

import (
	"github.com/open-cmi/cmmns/pkg/evchan"
	"github.com/open-cmi/cmmns/service/initial"
)

var echan *evchan.EventChan

func Register(event string, handler func(ev string, data interface{})) error {
	return echan.RegisterEvent(event, handler)
}

func Notify(event string, data interface{}) {
	echan.NotifyEvent(event, data)
}

func Init() error {
	echan.Run()
	return nil
}

func init() {
	echan = evchan.NewEventChan()
	initial.Register("chan-event", initial.DefaultPriority, Init)
}
