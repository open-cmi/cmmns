package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// SenderEmailInstance migrate
type SenderEmailInstance struct {
}

// Up up migrate
func (mi SenderEmailInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS sender_email (
			id varchar(64) not NULL primary key,
			server varchar(256) not NULL default '',
			port integer not NULL default 465,
			sender varchar(256) not NULL default '',
			password varchar(256) not null default '',
			use_tls boolean default false
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi SenderEmailInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS sender_email`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20231018161906",
		Description: "sender_email",
		Ext:         "go",
		Instance:    SenderEmailInstance{},
	})
}
