package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/cmmns/essential/migrate"
)

// LicenseAddUserInstance migrate
type LicenseAddUserInstance struct {
}

// Up up migrate
func (mi LicenseAddUserInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		ALTER TABLE license ADD COLUMN user varchar(100) default '';
	`
	db.Exec(sqlClause)
	return nil
}

// Down down migrate
func (mi LicenseAddUserInstance) Down(db *sqlx.DB) error {
	return nil
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:           "20250506154800",
		Description:   "license_add_user",
		Ext:           "go",
		Instance:      LicenseAddUserInstance{},
		Ignore:        false,
		AlterOpertion: true,
	})
}
