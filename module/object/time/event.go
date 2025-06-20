package time

import (
	"errors"
	"time"

	"github.com/open-cmi/cmmns/essential/logger"
)

var Registers map[string]func(string, bool) = make(map[string]func(string, bool))

func RegisterEvent(mod string, callback func(string, bool)) error {
	_, ok := Registers[mod]
	if ok {
		return errors.New("mod has been registered")
	}
	Registers[mod] = callback
	return nil
}

var allObjs map[string]*TimeObject = make(map[string]*TimeObject)

func AddObject(obj *TimeObject) {
	allObjs[obj.Name] = obj
}

func RemoveObject(obj *TimeObject) {
	delete(allObjs, obj.Name)
}

func CheckObjectStatus() {
	for _, obj := range allObjs {
		if obj.active != obj.IsActive() {
			logger.Infof("time object status change from %v to %v \n", obj.active, obj.IsActive())
			for _, cb := range Registers {
				cb(obj.Name, obj.IsActive())
				obj.active = obj.IsActive()
			}
		}
	}
}

var gStopChan chan bool = make(chan bool)

func ObjectEventLoop() {
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-t.C:
			CheckObjectStatus()
		case <-gStopChan:
			t.Stop()
			return
		}
	}
}
