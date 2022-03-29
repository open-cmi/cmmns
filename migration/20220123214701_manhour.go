package migration

import (
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// ManhourInstance migrate
type ManhourInstance struct {
}

// Up up migrate
func (mi ManhourInstance) Up() error {
	db := global.DB

	sqlClause := `
		CREATE TABLE IF NOT EXISTS manhour (
			id CHAR(64) NOT NULL PRIMARY KEY,
			date BIGINT NOT NULL DEFAULT 0,
			start_time BIGINT NOT NULL DEFAULT 0,
			end_time BIGINT NOT NULL DEFAULT 0,
			content TEXT NOT NULL DEFAULT '',
			created_time BIGINT NOT NULL DEFAULT 0,
			updated_time BIGINT NOT NULL DEFAULT 0
		)
	`
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi ManhourInstance) Down() error {
	db := global.DB

	sqlClause := `DROP TABLE IF EXISTS manhour`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20220123214701",
		Description: "manhour",
		Ext:         "go",
		Instance:    ManhourInstance{},
	})
}
