package config

import (
	"encoding/json"

	"github.com/open-cmi/cmmns/pkg/confparser"
)

var gConfCtx *confparser.Context

// Init config init
func Init(configfile string) error {

	err := gConfCtx.Load(configfile)
	return err
}

func Save() {
	if gConfCtx != nil {
		gConfCtx.Save()
	}
}

// RegisterConfig register config
func RegisterConfig(name string, parseFunc func(json.RawMessage) error, saveFunc func() json.RawMessage) error {

	if gConfCtx == nil {
		gConfCtx = confparser.NewContext()
	}

	var opt confparser.Option
	opt.Name = name
	opt.ParseFunc = parseFunc
	opt.SaveFunc = saveFunc
	return gConfCtx.Register(&opt)
}
