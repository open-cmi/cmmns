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

var gConf Config

// Init db init
func (c *Config) Init() error {
	var dbconf database.Config
	dbconf.Type = c.Type
	dbconf.File = c.File
	dbconf.Host = c.Host
	dbconf.Port = c.Port
	dbconf.User = c.User
	dbconf.Password = c.Password
	dbconf.Database = c.Database

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
	config.RegisterConfig("model", &gConf)
}
