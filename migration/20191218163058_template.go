package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// TemplateInstance migrate
type TemplateInstance struct {
}

// Up up migrate
func (mi TemplateInstance) Up(db *sqlx.DB) error {

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
func (mi TemplateInstance) Down(db *sqlx.DB) error {

	_, err := db.Exec(`
		DROP TABLE IF EXISTS template;
	`)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20191218163058",
		Description: "template",
		Ext:         "go",
		Instance:    TemplateInstance{},
	})
}
