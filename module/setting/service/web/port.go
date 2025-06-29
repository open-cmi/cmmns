package web

import (
	"github.com/open-cmi/cmmns/module/nginxconf"
	"github.com/open-cmi/gobase/essential/logger"
)

func SetServicePort(req *SetServicePortRequest) error {
	m := GetServicePortModel()
	if m == nil {
		m = NewServicePortModel()
	}
	err := nginxconf.ApplyServicePort(req.HTTPPort, req.HTTPSPort)
	if err != nil {
		logger.Errorf("apply port service failed: %s\n", err.Error())
		return err
	}
	m.HTTPPort = req.HTTPPort
	m.HTTPSPort = req.HTTPSPort
	err = m.Save()
	return err
}

func GetServicePort() *ServicePortModel {
	m := GetServicePortModel()
	return m
}
