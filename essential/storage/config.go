package storage

import (
	"github.com/open-cmi/cmmns/essential/storage/rdb"
	"github.com/open-cmi/cmmns/essential/storage/sqldb"
)

func Init() error {

	sqldb.Init()
	rdb.Init()
	return nil
}
