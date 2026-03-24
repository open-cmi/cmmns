package dbkv

import (
	"errors"
	"fmt"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/sqldb"
)

func Set(key string, value string) error {
	db := sqldb.GetDB()
	if db == nil {
		return errors.New("database not initialized")
	}

	sqlStatement := `
		INSERT INTO k_v_table (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE
		SET value = $2;
	`

	_, err := db.Exec(sqlStatement, key, value)
	if err != nil {
		return fmt.Errorf("failed to set key-value pair: %w", err)
	}

	return nil
}

func Get(key string) (string, bool) {
	db := sqldb.GetDB()
	if db == nil {
		return "", false
	}

	var value string
	queryClause := `SELECT value FROM k_v_table WHERE key=$1`
	row := db.QueryRow(queryClause, key)
	if row == nil {
		return "", false
	}
	err := row.Scan(&value)
	if err != nil {
		logger.Debugf("dbkv Get row scan failed: %s\n", err.Error())
		return "", false
	}
	return value, true
}
