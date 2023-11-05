package rbac

import (
	"errors"

	"github.com/open-cmi/cmmns/service/business"
)

var ModuleMapping map[string]string = make(map[string]string)

func RegisterModule(mod string, description string) error {
	_, ok := ModuleMapping[mod]
	if ok {
		return errors.New("module has been registered")
	}
	ModuleMapping[mod] = description
	return nil
}

func Init() error {
	for name, desc := range ModuleMapping {
		m := GetModule(name)
		if m == nil {
			m = NewModule(name, desc)
		}
		err := m.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	business.Register("rbac", Init)
}
