package sqldb

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/goutils/database"
	"github.com/open-cmi/goutils/database/dbsql"

	"github.com/jmoiron/sqlx"
)

// gConfDB sql db
var gConfDB *sqlx.DB
var gDataDB *sqlx.DB

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
var gDataModel Config

// GetDB get db
func GetConfDB() *sqlx.DB {
	return gConfDB
}

func GetDataDB() *sqlx.DB {
	return gDataDB
}

// Parse db init
func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConfModel)
	if err != nil {
		return err
	}
	var dbconf database.Config
	dbconf.Type = gConfModel.Type
	dbconf.File = gConfModel.File
	dbconf.Host = gConfModel.Host
	dbconf.Port = gConfModel.Port
	dbconf.User = gConfModel.User
	dbconf.Password = gConfModel.Password
	dbconf.Database = gConfModel.Database

	dbi, err := dbsql.SQLInit(&dbconf)
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

// Parse db init
func DataModelParse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gDataModel)
	if err != nil {
		return err
	}
	var dbconf database.Config
	dbconf.Type = gDataModel.Type
	dbconf.File = gDataModel.File
	dbconf.Host = gDataModel.Host
	dbconf.Port = gDataModel.Port
	dbconf.User = gDataModel.User
	dbconf.Password = gDataModel.Password
	dbconf.Database = gDataModel.Database

	dbi, err := dbsql.SQLInit(&dbconf)
	if err != nil {
		return err
	}
	gDataDB = dbi

	return nil
}

func DataModelSave() json.RawMessage {
	raw, _ := json.Marshal(&gDataModel)
	return raw
}

func init() {
	config.RegisterConfig("model", Parse, Save)
	config.RegisterConfig("data_model", DataModelParse, DataModelSave)
}
