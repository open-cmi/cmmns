package config

import (
	"encoding/json"
	"errors"

	"github.com/open-cmi/goutils/confparser"
)

type Feature interface {
	Init() error
}

// Config config
type Config struct {
	parser   *confparser.Parser
	Features map[string]Feature
}

var gConf *Config
var configMapping map[string]Feature = make(map[string]Feature)

// Init config init
func Init(configfile string) error {
	parser := confparser.New(configfile)
	if parser == nil {
		return errors.New("open file failed")
	}

	gConf = new(Config)
	gConf.parser = parser
	gConf.Features = configMapping

	var tmpConf map[string]json.RawMessage = make(map[string]json.RawMessage)
	err := gConf.parser.Load(&tmpConf)
	if err != nil {
		return err
	}

	for name, value := range tmpConf {
		moduleConfig, found := gConf.Features[name]
		if !found {
			continue
		}
		err := json.Unmarshal(value, moduleConfig)
		if err != nil {
			return err
		}
		moduleConfig.Init()
	}
	return nil
}

// Save save config
func Save() {
	gConf.parser.Save(gConf.Features)
}

// RegisterConfig register config
func RegisterConfig(name string, conf Feature) error {
	_, found := configMapping[name]
	if found {
		return errors.New("config " + name + " has been registered")
	}
	configMapping[name] = conf
	return nil
}
