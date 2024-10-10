package evchan

import (
	"errors"
	"fmt"
)

type EventData struct {
	Event string
	Data  interface{}
}

type EventChan struct {
	ExitChan chan int
	Chan     chan *EventData
	Handlers map[string]func(ev string, data interface{})
	Running  bool
}

func NewEventChan() *EventChan {
	return &EventChan{
		ExitChan: make(chan int, 1),
		Chan:     make(chan *EventData, 10),
		Handlers: make(map[string]func(ev string, data interface{})),
	}
}

func (ev *EventChan) NotifyEvent(event string, data interface{}) {
	ev.Chan <- &EventData{
		Event: event,
		Data:  data,
	}
}

func (ev *EventChan) RegisterEvent(event string, handler func(event string, data interface{})) error {
	_, ok := ev.Handlers[event]
	if ok {
		errmsg := fmt.Sprintf("event handler %s is existed", event)
		return errors.New(errmsg)
	}
	ev.Handlers[event] = handler
	return nil
}

func (ev *EventChan) Run() {
	ev.Running = true
	go func() {
		var loop bool = true
		for loop {
			select {
			case <-ev.ExitChan:
				loop = false
			case d := <-ev.Chan:
				event := d.Event
				handler, ok := ev.Handlers[event]
				if ok {
					handler(event, d.Data)
				}
			}
		}

		ev.Running = false
	}()
}
