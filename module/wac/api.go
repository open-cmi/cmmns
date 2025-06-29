package wac

import (
	"github.com/open-cmi/gobase/essential/logger"

	"github.com/open-cmi/cmmns/module/nginxconf"
)

func GetWAC() Model {
	m := Get()
	if m == nil {
		m = New()
		m.Mode = "blacklist"
	}
	return *m
}

func SetWAC(req *SetRequest) error {
	m := Get()
	if m == nil {
		m = New()
	}

	m.Enable = req.Enable
	m.Mode = req.Mode

	// 保存之前先应用nginx配置，nginx应用成功再保存
	var err error
	if req.Enable {
		err = nginxconf.ApplyNginxAccessControl(m.Mode)
	} else {
		err = nginxconf.CancelNginxAccessControl()
	}

	if err != nil {
		logger.Errorf("wac apply nginx access control failed: %s\n", err.Error())
		return err
	}

	err = m.Save()
	if err != nil {
		logger.Errorf("wac save failed: %s\n", err.Error())
		return err
	}

	return nil
}
