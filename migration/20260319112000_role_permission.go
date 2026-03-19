package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// RolePermissionInstance adds permission storage for API authorization.
type RolePermissionInstance struct{}

func (mi RolePermissionInstance) Up(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS role_permissions (
		role varchar(32) NOT NULL,
		perm varchar(256) NOT NULL,
		created_time integer NOT NULL default 0,
		PRIMARY KEY(role, perm)
	);
	`)
	if err != nil {
		return err
	}

	return nil
}

func (mi RolePermissionInstance) Down(db *sqlx.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS role_permissions;`)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20260319112000",
		Description: "role_permission",
		Ext:         "go",
		Instance:    RolePermissionInstance{},
	})
}
