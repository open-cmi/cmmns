package migration

import (
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// AuditLogInstance migrate
type AuditLogInstance struct {
	Name string
}

// Up up migrate
func (mi AuditLogInstance) Up() error {
	db := global.DB

	dbsql := `
		CREATE TABLE IF NOT EXISTS audit_log (
			id varchar(64) NOT NULL primary key,
			type integer NOT NULL default 0,
			ip varchar(64) NOT NULL default '',
			username varchar(100) NOT NULL default '',
			action text NOT NULL default '',
			timestamp integer NOT NULL default 0
		);
	`
	_, err := db.Exec(dbsql)
	return err
}

// Down down migrate
func (mi AuditLogInstance) Down() error {
	db := global.DB

	dbsql := `
		DROP TABLE IF EXISTS audit_log;
	`
	_, err := db.Exec(dbsql)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20211215135838",
		Description: "audit_log",
		Ext:         "go",
		Instance:    AuditLogInstance{},
	})
}
