package migration

import (
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"
)

// TemplateInstance migrate
type TemplateInstance struct {
}

// Up up migrate
func (mi TemplateInstance) Up() error {
	db := global.DB

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS template (
			id varchar(128) NOT NULL PRIMARY KEY default '',
			updated_time BIGINT NOT NULL default 0,
			created_time BIGINT NOT NULL default 0,
			name varchar(128) NOT NULL default ''
		);
	`)

	return err
}

// Down down migrate
func (mi TemplateInstance) Down() error {
	db := global.DB

	_, err := db.Exec(`
		DROP TABLE IF EXISTS template;
	`)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20191218163058",
		Description: "template",
		Ext:         "go",
		Instance:    TemplateInstance{},
	})
}
