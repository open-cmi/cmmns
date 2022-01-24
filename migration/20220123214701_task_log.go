package migration

import (
	"fmt"

	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// TaskLogInstance migrate
type TaskLogInstance struct {
}

// Up up migrate
func (mi TaskLogInstance) Up() error {
	db := global.DB

	sqlClause := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS task_log (
			id CHAR(64) NOT NULL PRIMARY KEY,
			type VARCHAR(32) NOT NULL DEFAULT '',
			content VARCHAR(512) NOT NULL DEFAULT '',
			file VARCHAR(512) NOT NULL DEFAULT '',
			ctime INTEGER NOT NULL DEFAULT 0
		)
	`)
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi TaskLogInstance) Down() error {
	db := global.DB

	sqlClause := fmt.Sprintf("DROP TABLE IF EXISTS task_log")
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20220123214701",
		Description: "task_log",
		Ext:         "go",
		Instance:    TaskLogInstance{},
	})
}
