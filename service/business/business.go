package business

import (
	"errors"
	"fmt"
	"sort"
)

const (
	DefaultPriority = 5
)

type Business struct {
	Init     func() error
	Priority int
	Name     string
}

var businesses []Business

func Init() error {
	var err error

	sort.SliceStable(businesses, func(i int, j int) bool {
		bz1 := businesses[i]
		bz2 := businesses[j]
		return bz1.Priority < bz2.Priority
	})

	for i := range businesses {
		bz := &businesses[i]
		err = bz.Init()
		if err != nil {
			errmsg := fmt.Sprintf("business %s init failed: %s", bz.Name, err.Error())
			return errors.New(errmsg)
		}
	}

	return nil
}

func Register(name string, priority int, fn func() error) error {
	for i := range businesses {
		bz := &businesses[i]
		if bz.Name == name {
			errmsg := fmt.Sprintf("business %s has been registered", name)
			return errors.New(errmsg)
		}
	}

	businesses = append(businesses, Business{
		Init:     fn,
		Name:     name,
		Priority: priority,
	})
	return nil
}
