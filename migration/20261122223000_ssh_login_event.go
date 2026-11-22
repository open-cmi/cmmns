package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// SSHLoginEventInstance migrate
type SSHLoginEventInstance struct{}

// Up up migrate
func (mi SSHLoginEventInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS ssh_login_event (
			id VARCHAR(64) NOT NULL PRIMARY KEY,
			timestamp INTEGER NOT NULL,
			username VARCHAR(64) NOT NULL DEFAULT '',
			rhost VARCHAR(64) NOT NULL DEFAULT '',
			port INTEGER NOT NULL DEFAULT 0,
			status VARCHAR(16) NOT NULL DEFAULT '',
			auth_type VARCHAR(32) NOT NULL DEFAULT '',
			is_invalid INTEGER NOT NULL DEFAULT 0,
			fail_reason VARCHAR(128) NOT NULL DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS idx_ssh_login_timestamp ON ssh_login_event(timestamp);
		CREATE INDEX IF NOT EXISTS idx_ssh_login_rhost ON ssh_login_event(rhost);
		CREATE INDEX IF NOT EXISTS idx_ssh_login_status ON ssh_login_event(status);
		CREATE INDEX IF NOT EXISTS idx_ssh_login_username ON ssh_login_event(username);
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi SSHLoginEventInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS ssh_login_event`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:           "20261122223000",
		Description:   "ssh_login_event",
		Ext:           "go",
		Instance:      SSHLoginEventInstance{},
		Ignore:        false,
		AlterOpertion: false,
	})
}
