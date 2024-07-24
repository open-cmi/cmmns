package hostname

import (
	"errors"
	"os/exec"

	"github.com/open-cmi/cmmns/essential/logger"
)

type SetRequest struct {
	Hostname string `json:"hostname"`
}

func Set(req *SetRequest) error {
	var args []string = []string{"set-hostname", req.Hostname}
	err := exec.Command("hostnamectl", args...).Run()
	if err != nil {
		logger.Errorf("set hostname failed: %s\n", err.Error())
		return errors.New("set hostname failed")
	}
	return nil
}
