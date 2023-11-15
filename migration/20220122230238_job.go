package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// JobInstance migrate
type JobInstance struct {
}

// Up up migrate
func (mi JobInstance) Up(db *sqlx.DB) error {

	sqlClause := `
		CREATE TABLE IF NOT EXISTS job (
			id char(64) NOT NULL PRIMARY KEY,
			type VARCHAR(32) NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			sched_group VARCHAR(64) NOT NULL DEFAULT '',
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
func (mi JobInstance) Down(db *sqlx.DB) error {

	sqlClause := `DROP TABLE IF EXISTS job`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20220122230238",
		Description: "job",
		Ext:         "go",
		Instance:    JobInstance{},
	})
}
