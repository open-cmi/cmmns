package migration

import (
	"fmt"

	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// SecretKeyInstance migrate
type SecretKeyInstance struct {
}

// Up up migrate
func (mi SecretKeyInstance) Up() error {
	db := global.DB

	sqlClause := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS secret_key (
			id char(64) NOT NULL PRIMARY KEY,
			name VARCHAR(256) NOT NULL unique DEFAULT '',
			user_id CHAR(64) NOT NULL DEFAULT '',
			key_type VARCHAR(12) NOT NULL DEFAULT '',
			key_length integer DEFAULT 2048,
			comment VARCHAR(128) NOT NULL DEFAULT '',
			passphrase VARCHAR(128) NOT NULL DEFAULT '',
			confirmation VARCHAR(128) NOT NULL DEFAULT '',
			private_key TEXT NOT NULL DEFAULT '',
			public_key TEXT NOT NULL DEFAULT '',
			ctime INT DEFAULT 0,
			utime INT DEFAULT 0
		)
	`)
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi SecretKeyInstance) Down() error {
	db := global.DB

	sqlClause := fmt.Sprintf(`
		DROP TABLE IF EXISTS secret_key
	`)
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20220121065958",
		Description: "secret_key",
		Ext:         "go",
		Instance:    SecretKeyInstance{},
	})
}
