package db

import (
	"github.com/open-cmi/cmmns/config"

	"github.com/open-cmi/goutils/database"
	"github.com/open-cmi/goutils/database/dbsql"

	"github.com/jmoiron/sqlx"
)

// DB sql db
var DB *sqlx.DB

// Init db init
func Init() error {
	var dbconf database.Config
	model := config.GetConfig().Model
	dbconf.Type = model.Type
	dbconf.File = model.File
	dbconf.Host = model.Host
	dbconf.Port = model.Port
	dbconf.User = model.User
	dbconf.Password = model.Password
	dbconf.Database = model.Database

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
