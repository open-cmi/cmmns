package service

import (
	"github.com/open-cmi/cmmns/module/nginxconf"
	"github.com/open-cmi/cmmns/service/initial"
)

func Init() error {
	m := GetServicePortModel()
	if m != nil {
		return nginxconf.ApplyServicePort(m.HTTPPort, m.HTTPSPort)
	}
	return nil
}

func init() {
	initial.Register("service-port-setting", initial.DefaultPriority, Init)
}
