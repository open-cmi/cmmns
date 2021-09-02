package db

import (
	"database/sql"
)

// DB sql db
var DB *sql.DB

// GetDB get db
func GetDB() *sql.DB {
	return DB
}
