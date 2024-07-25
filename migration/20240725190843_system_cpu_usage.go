package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// CPUUsageInstance migrate
type CPUUsageInstance struct {
}

// Up up migrate
func (mi CPUUsageInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS system_cpu_usage (
			dev_id varchar(128) NOT NULL,
			step int NOT NULL DEFAULT 0,
			updated_time bigint not null default 0,
			cpu_usage real default 0,
			primary key(dev_id,step)
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi CPUUsageInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS system_cpu_usage`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20240725190843",
		Description: "system_cpu_usage",
		Ext:         "go",
		Instance:    CPUUsageInstance{},
	})
}
