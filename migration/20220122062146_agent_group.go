package migration

import (
	"fmt"

	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// AgentGroupInstance migrate
type AgentGroupInstance struct {
}

// Up up migrate
func (mi AgentGroupInstance) Up() error {
	db := global.DB

	sqlClause := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS agent_group (
			id char(64) NOT NULL PRIMARY KEY,
			name VARCHAR(256) NOT NULL unique DEFAULT '',
			description VARCHAR(256) NOT NULL DEFAULT ''
		)
	`)
	_, err := db.Exec(sqlClause)
	return err
}

// Down down migrate
func (mi AgentGroupInstance) Down() error {
	db := global.DB

	dbsql := fmt.Sprintf(`DROP TABLE IF EXISTS agent_group`)
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
