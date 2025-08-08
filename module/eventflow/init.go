package eventflow

import (
	"github.com/open-cmi/gobase/initial"
)

func Init() error {
	if GetConf().Relay == "rdb" {
		err := RdbInit()
		if err != nil {
			return err
		}
	}

	for k, v := range GetConf().Channels {
		go EventFlowRecvProc(k, v)
	}

	return nil
}

func init() {
	initial.Register("event-flow", initial.PhaseDefault, Init)
}
