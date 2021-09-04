package db

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
)

// DB sql db
var DB *sql.DB

// Cache redis Cache
var Cache *redis.Client

// GetDB get db
func GetDB() *sql.DB {
	return DB
}

// GetCache get cache
func GetCache() *redis.Client {
	return Cache
}
