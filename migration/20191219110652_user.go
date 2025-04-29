package migration

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/cmmns/essential/migrate"

	"github.com/jameskeane/bcrypt"
)

// UserInstance migrate
type UserInstance struct {
}

// SyncAdminData sync data
func (mi UserInstance) SyncOperatorData(db *sqlx.DB) error {
	id := uuid.New().String()

	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash("123456", salt)
	itime := time.Now().Unix()
	dbsql := fmt.Sprintf(`
		INSERT INTO users (id, username, password, email, role, status, activate, itime, utime, description) 
			values ('%s', 'operator', '%s', 'admin@localhost',
			'operator', 'offline', true, %d, %d, 'operator');
  `, id, hash, itime, itime)
	_, err := db.Exec(dbsql)

	return err
}

// SyncAdminData sync data
func (mi UserInstance) SyncAdminData(db *sqlx.DB) error {
	id := "87fd7602-d20d-4ccb-80c5-b554cae79ce8"

	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash("admin12345678", salt)
	itime := time.Now().Unix()
	dbsql := fmt.Sprintf(`
		INSERT INTO users (id, username, password, email, role, status, activate, itime, utime, description) 
			values ('%s', 'admin', '%s', 'admin@localhost',
			'admin', 'offline', true, %d, %d, 'administrator');
  `, id, hash, itime, itime)
	_, err := db.Exec(dbsql)

	return err
}

// Up up migrate
func (mi UserInstance) Up(db *sqlx.DB) error {

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id varchar(36) primary key,
		username varchar(100) NOT NULL UNIQUE,
		password varchar(100) NOT NULL,
		email varchar(100) UNIQUE NOT NULL,
		role varchar(32) NOT NULL default '',
		status varchar(32) NOT NULL default 'offline',
		activate bool not NULL default false,
		itime integer NOT NULL default 0,
		utime integer NOT NULL default 0,
		password_change_time integer not NULL default 0,
		description text NOT NULL DEFAULT ''
      );
	`)
	if err == nil {
		errSync1 := mi.SyncAdminData(db)
		errSync2 := mi.SyncOperatorData(db)
		if errSync1 != nil {
			return errSync1
		}
		if errSync2 != nil {
			return errSync2
		}
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
