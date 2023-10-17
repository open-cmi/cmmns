package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// ManhourInstance migrate
type ManhourInstance struct {
}

// Up up migrate
func (mi ManhourInstance) Up(db *sqlx.DB) error {

	sqlClause := `
		CREATE TABLE IF NOT EXISTS manhour (
			id CHAR(64) NOT NULL PRIMARY KEY,
			date BIGINT NOT NULL DEFAULT 0,
			start_time BIGINT NOT NULL DEFAULT 0,
			end_time BIGINT NOT NULL DEFAULT 0,
			content TEXT NOT NULL DEFAULT '',
			created_time BIGINT NOT NULL DEFAULT 0,
			updated_time BIGINT NOT NULL DEFAULT 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi ManhourInstance) Down(db *sqlx.DB) error {

	sqlClause := `DROP TABLE IF EXISTS manhour`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20220123214701",
		Description: "manhour",
		Ext:         "go",
		Instance:    ManhourInstance{},
	})
}
