package migration

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/migrate"

	"github.com/jameskeane/bcrypt"
)

// UserInstance migrate
type UserInstance struct {
}

// SyncData sync data
func (mi UserInstance) SyncData(db *sqlx.DB) error {
	id := uuid.New().String()

	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash("admin12345678", salt)
	itime := time.Now().Unix()
	dbsql := fmt.Sprintf(`
		INSERT INTO users (id, username, password, email, role, status, itime, utime, description) 
			values ('%s', 'admin', '%s', 'admin@admin.com',
			'0', 1, %d, %d, '');
  `, id, hash, itime, itime)
	_, err := db.Exec(dbsql)

	return err
}

// Up up migrate
func (mi UserInstance) Up(db *sqlx.DB) error {

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id varchar(36) primary key,
		username varchar(100) NOT NULL,
		password varchar(200) NOT NULL,
		email varchar(100) UNIQUE NOT NULL,
		role integer NOT NULL default 100,
		status integer NOT NULL default 0,
		itime integer NOT NULL default 0,
		utime integer NOT NULL default 0,
		description text NOT NULL DEFAULT ''
      );
	`)
	if err == nil {
		mi.SyncData(db)
	}
	return err
}

// Down down migrate
func (mi UserInstance) Down(db *sqlx.DB) error {

	_, err := db.Exec(`
		DROP TABLE IF EXISTS users;
	`)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:         "20191219110652",
		Description: "user",
		Ext:         "go",
		Instance:    UserInstance{},
	})
}
