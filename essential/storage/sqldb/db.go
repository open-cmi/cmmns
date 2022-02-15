package sqldb

import (
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

var moduleConfig Config

// Init db init
func Init() error {
	var dbconf database.Config
	dbconf.Type = moduleConfig.Type
	dbconf.File = moduleConfig.File
	dbconf.Host = moduleConfig.Host
	dbconf.Port = moduleConfig.Port
	dbconf.User = moduleConfig.User
	dbconf.Password = moduleConfig.Password
	dbconf.Database = moduleConfig.Database

	dbi, err := dbsql.SQLInit(&dbconf)
	if err != nil {
		return err
	}
	DB = dbi

	return nil
}

// GetDB get db
func GetDB() *sqlx.DB {
	return DB
}

func init() {
	config.RegisterConfig("model", &moduleConfig)
}
