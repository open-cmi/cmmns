package sqldb

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/essential/config"
	"github.com/open-cmi/goutils/database"
	"github.com/open-cmi/goutils/database/dbsql"

	"github.com/jmoiron/sqlx"
)

// DB sql db
var DB *sqlx.DB

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

var gConf Config

// GetDB get db
func GetDB() *sqlx.DB {
	return DB
}

// Init db init
func Init(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}
	var dbconf database.Config
	dbconf.Type = gConf.Type
	dbconf.File = gConf.File
	dbconf.Host = gConf.Host
	dbconf.Port = gConf.Port
	dbconf.User = gConf.User
	dbconf.Password = gConf.Password
	dbconf.Database = gConf.Database

	dbi, err := dbsql.SQLInit(&dbconf)
	if err != nil {
		return err
	}
	DB = dbi

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	config.RegisterConfig("model", Init, Save)
}
