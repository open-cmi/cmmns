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
	return json.Unmarshal(mess, &gConf)
}

// Save marshals syslog configuration to json
func Save() json.RawMessage {
	v, _ := json.Marshal(gConf)
	return v
}

func init() {
	config.RegisterConfig("syslog", Parse, Save)
}
