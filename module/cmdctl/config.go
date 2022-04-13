package cmdctl

import (
	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/essential/logger"
)

var manager *Manager

// Config process config
type Config struct {
	Name       string `json:"name"`
	ExecStart  string `json:"exec_start"`
	RestartSec int    `json:"restart_sec"`
	StopSignal int    `json:"stop_signal"`
}

type CommandConfig struct {
	Services []Config `json:"services"`
}

func (cc *CommandConfig) Init() error {
	for _, s := range cc.Services {
		manager.AddProcess(&s)
		err := manager.StartProcess(s.Name)
		if err != nil {
			logger.Errorf("start process failed: %s\n", err.Error())
			return err
		}
	}
	return nil
}

var gConf CommandConfig

func init() {
	manager = NewManager()

	config.RegisterConfig("process", &gConf)
}
