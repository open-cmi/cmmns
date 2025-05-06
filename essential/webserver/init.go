package webserver

import (
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/initial"
)

func Init() error {
	// start web service
	s := New()

	// Init
	err := s.Init()
	if err != nil {
		logger.Errorf("web server init %s", err.Error())
		return err
	}

	// Run
	err = s.Run()
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

func init() {
	initial.Register("webserver", initial.PhaseWebserver, Init)
}
