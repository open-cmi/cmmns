package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// SyslogInstance migrate
type SyslogInstance struct {
}

// Up up migrate
func (mi SyslogInstance) Up(db *sqlx.DB) error {
	dbsql := `
		CREATE TABLE IF NOT EXISTS syslog (
			id varchar(64) NOT NULL primary key,
			log_time integer NOT NULL default 0,
			process varchar(64) NOT NULL default '',
			pid integer NOT NULL default 0,
			hostname varchar(64) NOT NULL default '',
			level varchar(16) NOT NULL default '',
			message text NOT NULL default '',
			created_at integer NOT NULL default 0
		);
		CREATE INDEX IF NOT EXISTS idx_syslog_process_time ON syslog(process, log_time);
		CREATE INDEX IF NOT EXISTS idx_syslog_created_at ON syslog(created_at);
	`
	_, err := db.Exec(dbsql)
	return err
}

// Down down migrate
func (mi SyslogInstance) Down(db *sqlx.DB) error {
	dbsql := `
		DROP TABLE IF EXISTS syslog;
	`
	_, err := db.Exec(dbsql)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20261122132000",
		Description: "syslog",
		Ext:         "go",
		Instance:    SyslogInstance{},
	})
}
