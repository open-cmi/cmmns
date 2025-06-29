package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// WhiteListInstance migrate
type WhiteListInstance struct {
}

// Up up migrate
func (mi WhiteListInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS wac_whitelist (
			address varchar(128) NOT NULL PRIMARY KEY,
			timestamp bigint not null default 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi WhiteListInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS wac_whitelist`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20240611164455",
		Description: "wac_whitelist",
		Ext:         "go",
		Instance:    WhiteListInstance{},
	})
}
