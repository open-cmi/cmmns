package storage

import (
	"github.com/open-cmi/cmmns/storage/db"
	"github.com/open-cmi/cmmns/storage/rdb"
)

func Init() {
	db.Init()
	rdb.Init()
}
