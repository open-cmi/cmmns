package initial

import (
	"errors"
	"fmt"
	"sort"
)

const (
	PhaseZero       = 0
	PhaseOne        = 1
	PhaseTwo        = 2
	DefaultPriority = 5
)

type Business struct {
	Init     func() error
	Priority int
	Name     string
}

var initiales []Business

func Init() error {
	var err error

	sort.SliceStable(initiales, func(i int, j int) bool {
		bz1 := initiales[i]
		bz2 := initiales[j]
		return bz1.Priority < bz2.Priority
	})

	for i := range initiales {
		bz := &initiales[i]
		err = bz.Init()
		if err != nil {
			errmsg := fmt.Sprintf("initial %s init failed: %s", bz.Name, err.Error())
			return errors.New(errmsg)
		}
	}

	return nil
}

func Register(name string, priority int, fn func() error) error {
	for i := range initiales {
		bz := &initiales[i]
		if bz.Name == name {
			errmsg := fmt.Sprintf("initial %s has been registered", name)
			return errors.New(errmsg)
		}
	}

	initiales = append(initiales, Business{
		Init:     fn,
		Name:     name,
		Priority: priority,
	})
	return nil
}
