package scheduler

import (
	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/essential/storage/rdb"
)

func init() {
	rdb.Register("scheduler", def.RDBScheduler)
}
