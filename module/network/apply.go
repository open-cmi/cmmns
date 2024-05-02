package network

import (
	"errors"
	"fmt"
)

func NetworkApply() error {
	// 这里要校验格式

	if gConf.Engine == "netplan" {
		return NetplanApply()
	} else if gConf.Engine == "networking" {
		return NetworkingApply()
	}
	errmsg := fmt.Sprintf("engine %s is not supported", gConf.Engine)
	return errors.New(errmsg)
}
