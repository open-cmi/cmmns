package business

import (
	"errors"
	"fmt"
)

var businesses map[string]func() error

func Init() error {
	var err error
	for name, fn := range businesses {
		err = fn()
		if err != nil {
			errmsg := fmt.Sprintf("business %s init failed: %s", name, err.Error())
			return errors.New(errmsg)
		}
	}
	return nil
}

func Register(name string, fn func() error) error {
	_, ok := businesses[name]
	if ok {
		errmsg := fmt.Sprintf("business %s has been registered", name)
		return errors.New(errmsg)
	}
	businesses[name] = fn
	return nil
}

func init() {
	businesses = make(map[string]func() error)
}
