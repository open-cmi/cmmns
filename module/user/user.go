package user

import (
	"github.com/open-cmi/cmmns/essential/rdb"
)

func init() {
	rdb.Register("user", rdb.RDBUser)
}
