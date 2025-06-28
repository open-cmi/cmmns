package network

import (
	"errors"
	"fmt"
)

func NetworkApply(conf *Config) error {
	// 这里要校验格式

	switch conf.Engine {
	case "netplan":
		return NetplanApply()
	case "networking":
		return NetworkingApply()
	default:
	}
	errmsg := fmt.Sprintf("engine %s is not supported", conf.Engine)
	return errors.New(errmsg)
}
