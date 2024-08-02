package events

import (
	"github.com/open-cmi/cmmns/pkg/evchan"
	"github.com/open-cmi/cmmns/service/initial"
)

var transfer *evchan.EventChan

func Register(event string, handler func(ev string, data interface{})) error {
	return transfer.RegisterEvent(event, handler)
}

func Notify(event string, data interface{}) {
	transfer.NotifyEvent(event, data)
}

func Init() error {
	transfer.Run()
	return nil
}

func init() {
	transfer = evchan.NewEventChan()
	initial.Register("events", initial.DefaultPriority, Init)
}
