package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/cmmns/essential/migrate"
)

// LicenseInstance migrate
type LicenseInstance struct {
}

// Up up migrate
func (mi LicenseInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS license (
			id char(64) NOT NULL PRIMARY KEY,
			customer VARCHAR(128) NOT NULL DEFAULT '',
			prod varchar(64) not NULL default '',
			version varchar(64) not null default '',
			modules text not NULL default '',
			model varchar(64) default '',
			expire_time BIGINT NOT NULL DEFAULT 0,
			mcode varchar(256) not null default '',
			created_time BIGINT NOT NULL DEFAULT 0,
			updated_time BIGINT NOT NULL DEFAULT 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi LicenseInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS license`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20240508050246",
		Description: "license",
		Ext:         "go",
		Instance:    LicenseInstance{},
	})
}
