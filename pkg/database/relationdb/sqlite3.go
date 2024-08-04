package relationdb

import (
	"os"
	"path/filepath"

	"github.com/open-cmi/cmmns/pkg/pathutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// SQLite3Init init
func SQLite3Init(conf *Config) (db *sqlx.DB, err error) {
	dbfile := conf.File
	if !filepath.IsAbs(conf.File) {
		dbfile = filepath.Join(pathutil.GetRootPath(), "data", conf.File)
	}

	// if filename is absolute path, use file name directly

	var file *os.File
	if !pathutil.IsExist(dbfile) {
		file, err = os.OpenFile(dbfile, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return
		}
		defer file.Close()
	}

	db, err = sqlx.Open("sqlite3", dbfile)
	return db, err
}
