package migration

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/gobase/essential/migrate"
)

// RBACInstance migrate
type RBACInstance struct {
}

// SyncData sync data
func (mi RBACInstance) SyncData(db *sqlx.DB) error {
	id := uuid.New().String()

	now := time.Now().Unix()
	// 添加admin用户
	dbsql := fmt.Sprintf(`
		INSERT INTO roles (id, name, created_time, updated_time, permisions, description) 
			values ('%s', 'admin', %d, %d, '*', 'administrators');`, id, now, now)
	_, err := db.Exec(dbsql)
	if err != nil {
		return err
	}

	// 添加operator角色
	id = uuid.New().String()
	dbsql = fmt.Sprintf(`
		INSERT INTO roles (id, name, created_time, updated_time, permisions, description) 
			values ('%s', 'operator', %d, %d, '*', 'operator');`, id, now, now)
	_, err = db.Exec(dbsql)
	if err != nil {
		return err
	}

	// 添加auditor角色
	id = uuid.New().String()
	dbsql = fmt.Sprintf(`
		INSERT INTO roles (id, name, created_time, updated_time, permisions, description) 
			values ('%s', 'auditor', %d, %d, '*', 'auditor');`, id, now, now)
	_, err = db.Exec(dbsql)
	if err != nil {
		return err
	}
	return nil
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

	return nil
}

// Down down migrate
func (mi RBACInstance) Down(db *sqlx.DB) error {

	_, err := db.Exec(`
		DROP TABLE IF EXISTS roles;
	`)

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
