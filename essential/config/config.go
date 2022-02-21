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

var config *Config
var defaultConfig map[string]Feature = make(map[string]Feature)

// Init config init
func Init(configfile string) error {
	parser := confparser.New(configfile)
	if parser == nil {
		return errors.New("open file failed")
	}

	config = new(Config)
	config.parser = parser
	config.Features = defaultConfig

	var tmpConf map[string]json.RawMessage = make(map[string]json.RawMessage)
	err := config.parser.Load(&tmpConf)
	if err != nil {
		return err
	}

	for name, value := range tmpConf {
		moduleConfig, found := config.Features[name]
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
	config.parser.Save(config.Features)
}

func RegisterConfig(name string, conf Feature) error {
	_, found := defaultConfig[name]
	if found {
		return errors.New("config " + name + " has been registered")
	}
	defaultConfig[name] = conf
	return nil
}
