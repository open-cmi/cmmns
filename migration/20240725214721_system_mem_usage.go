package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// MemUsageInstance migrate
type MemUsageInstance struct {
}

// Up up migrate
func (mi MemUsageInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS system_mem_usage (
			dev_id varchar(128) NOT NULL,
			step int NOT NULL DEFAULT 0,
			updated_time bigint not null default 0,
			mem_usage real default 0,
			primary key(dev_id,step)
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi MemUsageInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS system_mem_usage`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20240725214721",
		Description: "system_mem_usage",
		Ext:         "go",
		Instance:    MemUsageInstance{},
	})
}
