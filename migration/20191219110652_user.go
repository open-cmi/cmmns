package migration

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/open-cmi/migrate"
	"github.com/open-cmi/migrate/cmdopt"
	"github.com/open-cmi/migrate/global"

	"github.com/jameskeane/bcrypt"
)

// UserInstance migrate
type UserInstance struct {
}

// SyncData sync data
func (mi UserInstance) SyncData() error {
	db := global.DB
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
func (mi UserInstance) Up() error {
	db := global.DB

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
		err = mi.SyncData()
	}
	return err
}

// Down down migrate
func (mi UserInstance) Down() error {
	db := global.DB

	_, err := db.Exec(`
		DROP TABLE IF EXISTS users;
	`)
	return err
}

func init() {
	migrate.Register(&cmdopt.SeqInfo{
		Seq:         "20191219110652",
		Description: "user",
		Ext:         "go",
		Instance:    UserInstance{},
	})
}
