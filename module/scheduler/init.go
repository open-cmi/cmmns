package scheduler

import (
	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/essential/rdb"
)

func init() {
	rdb.Register("scheduler", def.RDBScheduler)
}
