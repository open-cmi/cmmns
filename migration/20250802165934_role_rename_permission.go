package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// RoleRenamePermissionInstance migrate
type RoleRenamePermissionInstance struct {
}

// Up up migrate
func (mi RoleRenamePermissionInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		ALTER TABLE roles rename COLUMN permisions to permissions;
	`
	db.Exec(sqlClause)
	return nil
}

// Down down migrate
func (mi RoleRenamePermissionInstance) Down(db *sqlx.DB) error {
	return nil
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:           "20250802165934",
		Description:   "role_rename_permission",
		Ext:           "go",
		Instance:      RoleRenamePermissionInstance{},
		Ignore:        false,
		AlterOpertion: true,
	})
}
