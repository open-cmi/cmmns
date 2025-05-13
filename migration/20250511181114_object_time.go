package migration

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/open-cmi/cmmns/essential/migrate"
	tmobj "github.com/open-cmi/cmmns/module/object/time"
)

// ObjectTimeInstance migrate
type ObjectTimeInstance struct {
}

func (mi ObjectTimeInstance) SyncData(db *sqlx.DB) error {
	end := time.Date(3999, 0, 0, 0, 0, 0, 0, time.UTC).Unix()

	var obj tmobj.AbsoluteTimeObject
	obj.Name = "always"
	obj.Description = "always time"
	obj.TimestampStart = 0
	obj.TimestampEnd = end
	v, _ := json.Marshal(obj)

	insertClause := fmt.Sprintf(`insert into object_time (name, refcnt, description, time_type, value) 
	values('%s', %d, '%s', %d, '%s')`, "always", 0, "always time", tmobj.TimeTypeAbsolute, v)
	_, err := db.Exec(insertClause)
	return err
}

// Up up migrate
func (mi ObjectTimeInstance) Up(db *sqlx.DB) error {
	sqlClause := `
		CREATE TABLE IF NOT EXISTS object_time (
			name VARCHAR(256) NOT NULL primary key,
			refcnt integer not null default 0,
			description varchar(256) not null default '',
			time_type integer default 0,
			value text default ''
		)
	`
	_, err := db.Exec(sqlClause)
	if err == nil {
		err = mi.SyncData(db)
	}
	return err
}

// Down down migrate
func (mi ObjectTimeInstance) Down(db *sqlx.DB) error {
	sqlClause := `DROP TABLE IF EXISTS object_time`
	_, err := db.Exec(sqlClause)
	return err
}

func init() {
	migrate.Register(&migrate.SeqInfo{
		Seq:           "20250511181114",
		Description:   "object_time",
		Ext:           "go",
		Instance:      ObjectTimeInstance{},
		Ignore:        false,
		AlterOpertion: false,
	})
}
