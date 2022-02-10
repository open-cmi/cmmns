package migration

import (
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// RealTimeTaskInstance migrate
type RealTimeTaskInstance struct {
}

// Up up migrate
func (mi RealTimeTaskInstance) Up() error {
	db := global.DB

	// result 0表示失败，1表示成功，2表示有成功有失败
	sqlClause := `
		CREATE TABLE IF NOT EXISTS realtime_task (
			id char(64) NOT NULL PRIMARY KEY,
			type VARCHAR(32) NOT NULL DEFAULT '',
			total INTEGER NOT NULL DEFAULT 0,
			success INTEGER NOT NULL DEFAULT 0,
			failed INTEGER NOT NULL DEFAULT 0,
			start_time  BIGINT NOT NULL DEFAULT 0,
			end_time BIGINT NOT NULL DEFAULT 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi RealTimeTaskInstance) Down() error {
	db := global.DB

	sqlClause := `DROP TABLE IF EXISTS realtime_task`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20220122230238",
		Description: "realtime_task",
		Ext:         "go",
		Instance:    RealTimeTaskInstance{},
	})
}
