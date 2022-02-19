package migration

import (
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// JobInstance migrate
type JobInstance struct {
}

// Up up migrate
func (mi JobInstance) Up() error {
	db := global.DB

	sqlClause := `
		CREATE TABLE IF NOT EXISTS job (
			id char(64) NOT NULL PRIMARY KEY,
			cron_id VARCHAR(128) NOT NULL DEFAULT '',
			type VARCHAR(32) NOT NULL DEFAULT '',
			priority INTEGER NOT NULL DEFAULT 0,
			content TEXT NOT NULL DEFAULT '',
			run_type VARCHAR(32) NOT NULL DEFAULT '',
			run_spec VARCHAR(64) NOT NULL DEFAULT '',
			sched_group VARCHAR(64) NOT NULL DEFAULT '',
			sched_object VARCHAR(64) NOT NULL DEFAULT '',
			state VARCHAR(32) NOT NULL DEFAULT '',
			count INTEGER NOT NULL DEFAULT 0,
			done INTEGER NOT NULL DEFAULT 0,
			code INTEGER NOT NULL DEFAULT 0,
			msg VARCHAR(512) NOT NULL DEFAULT '',
			result TEXT NOT NULL DEFAULT '',
			started_time  BIGINT NOT NULL DEFAULT 0,
			stopped_time BIGINT NOT NULL DEFAULT 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi JobInstance) Down() error {
	db := global.DB

	sqlClause := `DROP TABLE IF EXISTS job`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20220122230238",
		Description: "job",
		Ext:         "go",
		Instance:    JobInstance{},
	})
}
