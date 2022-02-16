package scheduler

import "github.com/open-cmi/cmmns/essential/storage/rdb"

func init() {
	rdb.Register("scheduler", 4)
}
