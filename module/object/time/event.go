package time

import (
	"errors"

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

func CheckObjectStatus(name string, data interface{}) {
	results, err := TimeObjectList()
	if err != nil {
		logger.Errorf("ticker run check status failed: %s\n", err.Error())
		return
	}
	for _, obj := range results {
		c, ok := allObjs[obj.Name]
		if !ok {
			allObjs[obj.Name] = &obj
			continue
		}
		if c.active != obj.IsActive() {
			for name, cb := range Registers {
				logger.Infof("start to run time object callback: %s\n", name)
				cb(obj.Name, obj.IsActive())
				c.active = obj.IsActive()
			}
		}
	}
}

// func init() {
// 	ticker.Register("time-object", "*/20 * * * * *", CheckObjectStatus, nil)
// }
