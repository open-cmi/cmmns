package pubsub

import (
	evbus "github.com/asaskevich/EventBus"
)

var bus evbus.Bus

func Subscribe(topic string, fn interface{}) {
	bus.Subscribe(topic, fn)
}

func Publish(topic string, args ...interface{}) {
	bus.Publish(topic, args...)
}

func init() {
	bus = evbus.New()
}
