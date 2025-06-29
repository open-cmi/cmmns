package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// KVInstance migrate
type KVInstance struct {
}

// Up up migrate
func (mi KVInstance) Up(db *sqlx.DB) error {

	sqlClause := `
		CREATE TABLE IF NOT EXISTS k_v_table (
			key VARCHAR(256) NOT NULL primary key,
			value text NOT NULL DEFAULT ''
		)
	`
	_, err := db.Exec(sqlClause)
	if err != nil {
		return err
	}

	dbsql := `CREATE INDEX kv_key_index ON k_v_table(key);`
	_, err = db.Exec(dbsql)
	if err != nil {
		return err
	}

	return err
}

// Down down migrate
func (mi KVInstance) Down(db *sqlx.DB) error {

	sqlClause := `DROP INDEX IF EXISTS kv_key_index`
	_, err := db.Exec(sqlClause)
	if err == nil {
		sqlClause := `DROP TABLE IF EXISTS k_v_table`
		_, err = db.Exec(sqlClause)
	}
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20230614144024",
		Description: "k_v_table",
		Ext:         "go",
		Instance:    KVInstance{},
	})
}
