package syslog

import (
	"encoding/json"

	"github.com/open-cmi/gobase/essential/config"
)

var gConf = Config{
	FilePath:  "",
	Processes: []string{},
	KeepDays:  30,
}

// Config syslog configuration
type Config struct {
	FilePath  string   `json:"file_path"`
	Processes []string `json:"processes"`
	KeepDays  int      `json:"keep_days"`
}

// Parse parses syslog configuration from json
func Parse(mess json.RawMessage) error {
	err := json.Unmarshal(mess, &gConf)
	if err != nil {
		return err
	}
	if gConf.KeepDays == 0 {
		gConf.KeepDays = 30
	}
	return nil
}

// Save marshals syslog configuration to json
func Save() json.RawMessage {
	v, _ := json.Marshal(gConf)
	return v
}

func init() {
	config.RegisterConfig("syslog", Parse, Save)
}
