package relationdb

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var maxConnections int = 50

type Config struct {
	Type     string
	File     string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// SQLInit init
func SQLInit(conf *Config) (db *sqlx.DB, err error) {
	// Mongo, _ = database.NewMongoDB("usual", 27017, "test")
	dbtype := conf.Type
	if dbtype == "sqlite3" {
		return SQLite3Init(conf)
	} else if dbtype == "postgresql" || dbtype == "pg" {
		return PostgresqlInit(conf)
	}
	return nil, errors.New("db not support")
}
