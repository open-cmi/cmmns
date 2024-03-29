package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// SecretKeyInstance migrate
type SecretKeyInstance struct {
}

// Up up migrate
func (mi SecretKeyInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS secret_key (
			id char(64) NOT NULL PRIMARY KEY,
			name VARCHAR(256) NOT NULL unique DEFAULT '',
			key_type VARCHAR(12) NOT NULL DEFAULT '',
			key_length integer DEFAULT 2048,
			comment VARCHAR(128) NOT NULL DEFAULT '',
			passphrase VARCHAR(128) NOT NULL DEFAULT '',
			confirmation VARCHAR(128) NOT NULL DEFAULT '',
			private_key TEXT NOT NULL DEFAULT '',
			public_key TEXT NOT NULL DEFAULT '',
			created_time INTEGER DEFAULT 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi SecretKeyInstance) Down(db *sqlx.DB) error {
	sqlClause := `
		DROP TABLE IF EXISTS secret_key
	`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20220121065958",
		Description: "secret_key",
		Ext:         "go",
		Instance:    SecretKeyInstance{},
	})
}
