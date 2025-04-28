package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/cmmns/essential/migrate"
)

// LicenseAddModelInstance migrate
type LicenseAddModelInstance struct {
}

// Up up migrate
func (mi LicenseAddModelInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		ALTER TABLE license ADD COLUMN model varchar(64) default '';
	`
	db.Exec(sqlClause)
	return nil
}

// Down down migrate
func (mi LicenseAddModelInstance) Down(db *sqlx.DB) error {
	return nil
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:           "20250427183410",
		Description:   "license_add_model",
		Ext:           "go",
		Instance:      LicenseAddModelInstance{},
		Ignore:        false,
		AlterOpertion: true,
	})
}
