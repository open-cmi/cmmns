package migration

import (
	"github.com/google/uuid"
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// AgentGroupInstance migrate
type AgentGroupInstance struct {
}

// SyncData sync data
func (agi AgentGroupInstance) SyncData() error {
	db := global.DB
	id := uuid.New().String()

	dbsql := `
		INSERT INTO agent_group (id, name, description) 
			values ($1, $2, $3);
  `
	_, err := db.Exec(dbsql, id, "default", "")

	return err
}

// Up up migrate
func (mi AgentGroupInstance) Up() error {
	db := global.DB

	sqlClause := `
		CREATE TABLE IF NOT EXISTS agent_group (
			id char(64) NOT NULL PRIMARY KEY,
			name VARCHAR(256) NOT NULL unique DEFAULT '',
			description VARCHAR(256) NOT NULL DEFAULT ''
		)
	`
	_, err := db.Exec(sqlClause)
	if err == nil {
		mi.SyncData()
	}
	return err
}

// Down down migrate
func (mi AgentGroupInstance) Down() error {
	db := global.DB

	dbsql := `DROP TABLE IF EXISTS agent_group`
	_, err := db.Exec(dbsql)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20220122062146",
		Description: "agent_group",
		Ext:         "go",
		Instance:    AgentGroupInstance{},
	})
}
