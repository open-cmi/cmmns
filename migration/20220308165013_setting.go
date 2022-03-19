package migration

import (
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// SettingInstance migrate
type SettingInstance struct {
}

// Up up migrate
func (mi SettingInstance) Up() error {
	db := global.DB

	sqlClause := `
		CREATE TABLE IF NOT EXISTS setting (
			id char(64) NOT NULL PRIMARY KEY,
			scope VARCHAR(32) NOT NULL DEFAULT '',
			belong VARCHAR(1024) NOT NULL DEFAULT '',
			key VARCHAR(256) NOT NULL DEFAULT '',
			value TEXT NOT NULL DEFAULT '',
			cfgseq int not NULL DEFAULT 0,
			created_time BIGINT NOT NULL default 0,
			updated_time BIGINT NOT NULL default 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi SettingInstance) Down() error {
	db := global.DB

	sqlClause := `DROP TABLE IF EXISTS setting`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20220308165013",
		Description: "setting",
		Ext:         "go",
		Instance:    SettingInstance{},
	})
}
