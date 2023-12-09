package sqldb

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/cmmns/pkg/database/relationdb"

	"github.com/jmoiron/sqlx"
)

// gConfDB sql db
var gConfDB *sqlx.DB

// DBConfig database model
type Config struct {
	Type     string `json:"type"`
	File     string `json:"file,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

var gConfModel Config

// GetDB get db
func GetConfDB() *sqlx.DB {
	return gConfDB
}

// Parse db init
func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConfModel)
	if err != nil {
		return err
	}
	var dbconf relationdb.Config
	dbconf.Type = gConfModel.Type
	dbconf.File = gConfModel.File
	dbconf.Host = gConfModel.Host
	dbconf.Port = gConfModel.Port
	dbconf.User = gConfModel.User
	dbconf.Password = gConfModel.Password
	dbconf.Database = gConfModel.Database

	dbi, err := relationdb.SQLInit(&dbconf)
	if err != nil {
		return err
	}
	gConfDB = dbi

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConfModel)
	return raw
}

func init() {
	config.RegisterConfig("model", Parse, Save)
}
