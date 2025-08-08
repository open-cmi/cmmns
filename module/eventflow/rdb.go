package eventflow

import (
	"errors"

	"github.com/open-cmi/gobase/essential/rdb"
)

var rcli *rdb.Client

func RdbInit() error {
	rcli = rdb.GetClient(0)
	if rcli == nil {
		return errors.New("rdb is not available")
	}
	return nil
}
