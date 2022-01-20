package migration

import (
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// SystemInfoMigration migrate
type SystemInfoMigration struct {
}

// Up up migrate
func (mi SystemInfoMigration) Up() error {
	db := global.DB

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS system_status (
			id CHAR(64) primary key,
			utime integer NOT NULL default 0,
			deviceid varchar(128) NOT NULL default '',
			is_master boolean default false,
			cpu_usage REAL DEFAULT 0,
			disk_total BIGINT DEFAULT 0,
			disk_used BIGINT DEFAULT 0,
			disk_used_percent REAL DEFAULT 0,
			mem_total BIGINT DEFAULT 0,
			mem_used BIGINT DEFAULT 0,
			mem_used_percent REAL DEFAULT 0,
			load_avg_1 REAL DEFAULT 0,
			load_avg_5 REAL DEFAULT 0,
			load_avg_15 REAL DEFAULT 0,
			net_sent BIGINT DEFAULT 0,
			net_recv BIGINT DEFAULT 0
		);
	`)

	return err
}

// Down down migrate
func (mi SystemInfoMigration) Down() error {
	db := global.DB

	_, err := db.Exec(`
		DROP TABLE IF EXISTS system_status;
	`)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20201226163058",
		Description: "create_system_status",
		Ext:         "go",
		Instance:    SystemInfoMigration{},
	})
}
