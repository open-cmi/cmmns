package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/cmmns/essential/migrate"
)

// BlackListInstance migrate
type BlackListInstance struct {
}

// Up up migrate
func (mi BlackListInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS wac_blacklist (
			address varchar(128) NOT NULL PRIMARY KEY,
			timestamp bigint not null default 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi BlackListInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS wac_blacklist`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20240611164511",
		Description: "wac_blacklist",
		Ext:         "go",
		Instance:    BlackListInstance{},
	})
}
