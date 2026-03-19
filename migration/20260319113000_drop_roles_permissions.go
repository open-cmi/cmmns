package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// DropRolesPermissionsColumnInstance drops legacy roles.permissions/permisions column.
type DropRolesPermissionsColumnInstance struct{}

func (mi DropRolesPermissionsColumnInstance) Up(db *sqlx.DB) error {
	// Make it safe across different historical schemas.
	_, _ = db.Exec(`ALTER TABLE roles DROP COLUMN IF EXISTS permissions;`)
	_, _ = db.Exec(`ALTER TABLE roles DROP COLUMN IF EXISTS permisions;`)
	return nil
}

func (mi DropRolesPermissionsColumnInstance) Down(db *sqlx.DB) error {
	// Not re-adding legacy column.
	return nil
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:           "20260319113000",
		Description:   "drop_roles_permissions",
		Ext:           "go",
		Instance:      DropRolesPermissionsColumnInstance{},
		AlterOpertion: true,
	})
}
