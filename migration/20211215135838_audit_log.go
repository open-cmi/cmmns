package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// AuditLogInstance migrate
type AuditLogInstance struct {
}

// Up up migrate
func (mi AuditLogInstance) Up(db *sqlx.DB) error {
	dbsql := `
		CREATE TABLE IF NOT EXISTS audit_log (
			id varchar(64) NOT NULL primary key,
			type integer NOT NULL default 0,
			ip varchar(64) NOT NULL default '',
			username varchar(100) NOT NULL default '',
			action text NOT NULL default '',
			result varchar(16) not null default '',
			timestamp integer NOT NULL default 0
		);
	`
	_, err := db.Exec(dbsql)
	return err
}

// Down down migrate
func (mi AuditLogInstance) Down(db *sqlx.DB) error {

	dbsql := `
		DROP TABLE IF EXISTS audit_log;
	`
	_, err := db.Exec(dbsql)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20211215135838",
		Description: "audit_log",
		Ext:         "go",
		Instance:    AuditLogInstance{},
	})
}
