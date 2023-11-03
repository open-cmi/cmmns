package migration

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"
)

// RBACInstance migrate
type RBACInstance struct {
}

// SyncData sync data
func (mi RBACInstance) SyncData(db *sqlx.DB) error {
	id := uuid.New().String()

	now := time.Now().Unix()
	dbsql := fmt.Sprintf(`
		INSERT INTO roles (id, name, created_time, updated_time, permisions, description) 
			values ('%s', 'admin', %d, %d, '*', 'administrators');
  `, id, now, now)
	_, err := db.Exec(dbsql)

	return err
}

// Up up migrate
func (mi RBACInstance) Up(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS roles (
		id char(36) primary key,
		name varchar(32) NOT NULL UNIQUE,
		created_time integer NOT NULL default 0,
		updated_time integer NOT NULL default 0,
		description text NOT NULL DEFAULT '',
		permisions text NOT NULL default '*'
      );
	`)
	if err == nil {
		err = mi.SyncData(db)
	}

	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS modules (
		name varchar(32) NOT NULL UNIQUE,
		description text NOT NULL DEFAULT ''
      );
	`)
	return err
}

// Down down migrate
func (mi RBACInstance) Down(db *sqlx.DB) error {

	_, err := db.Exec(`
		DROP TABLE IF EXISTS roles;
	`)

	if err == nil {
		_, err = db.Exec(`
			DROP TABLE IF EXISTS modules;
		`)
	}
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20191219100521",
		Description: "rbac",
		Ext:         "go",
		Instance:    RBACInstance{},
	})
}
