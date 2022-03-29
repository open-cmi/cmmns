package user

import (
	"github.com/open-cmi/cmmns/common/def"
	"github.com/open-cmi/cmmns/essential/rdb"
)

func init() {

	rdb.Register("user", def.RDBUser)
}
