package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// TokenRecordInstance migrate
type TokenRecordInstance struct {
}

// Up up migrate
func (mi TokenRecordInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS token_record (
			name varchar(256) NOT NULL PRIMARY KEY,
			expire_day integer NOT NULL DEFAULT 30,
			created_time integer not NULL default 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi TokenRecordInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS token_record`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20231117210145",
		Description: "token_record",
		Ext:         "go",
		Instance:    TokenRecordInstance{},
	})
}
