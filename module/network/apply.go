package network

import (
	"errors"
	"fmt"
)

func NetworkApply(conf *Config) error {
	// 这里要校验格式

	if conf.Engine == "netplan" {
		return NetplanApply(conf)
	}
	errmsg := fmt.Sprintf("engine %s is not supported", conf.Engine)
	return errors.New(errmsg)
}
