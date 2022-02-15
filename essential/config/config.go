package config

import (
	"encoding/json"
	"errors"

	"github.com/open-cmi/goutils/confparser"
)

// Config config
type Config struct {
	parser   *confparser.Parser
	Features map[string]interface{}
}

var config *Config

// Init config init
func Init(configfile string) error {
	parser := confparser.New(configfile)
	if parser == nil {
		return errors.New("open file failed")
	}

	conf := new(Config)
	conf.parser = parser
	conf.Features = defaultConfig

	var tmpConf map[string]json.RawMessage = make(map[string]json.RawMessage)
	err := conf.parser.Load(&tmpConf)
	if err != nil {
		return err
	}

	for name, value := range tmpConf {
		moduleConfig, found := conf.Features[name]
		if !found {
			continue
		}
		err := json.Unmarshal(value, moduleConfig)
		if err != nil {
			return err
		}
	}
	config = conf
	return nil
}

// Save save config
func Save() {
	config.parser.Save(config.Features)
}

var defaultConfig map[string]interface{} = make(map[string]interface{})

func RegisterConfig(name string, conf interface{}) error {
	_, found := defaultConfig[name]
	if found {
		return errors.New("config " + name + " has been registered")
	}
	defaultConfig[name] = conf
	return nil
}
