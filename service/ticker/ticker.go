package ticker

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
)

var initialized bool

type Ticker struct {
	Name string
	Spec string
	Func func()
}

var tickers map[string]Ticker = make(map[string]Ticker)
var cronMap map[string]*cron.Cron = make(map[string]*cron.Cron)

func Register(name string, spec string, f func()) error {
	_, found := tickers[name]
	if found {
		errMsg := fmt.Sprintf("ticker %s registered failed", name)
		return errors.New(errMsg)
	}
	tickers[name] = Ticker{
		Name: name,
		Spec: spec,
		Func: f,
	}
	if initialized {
		// 如果已经初始化了，此时需要立即创建
		ins := cron.New(cron.WithSeconds())
		ins.AddFunc(spec, f)
		ins.Start()
		cronMap[name] = ins
	}
	return nil
}

var cronInstance *cron.Cron

func Init() error {

	for _, ticker := range tickers {
		cronInstance = cron.New(cron.WithSeconds())
		cronInstance.AddFunc(ticker.Spec, ticker.Func)
		cronInstance.Start()
		cronMap[ticker.Name] = cronInstance
	}

	initialized = true
	return nil
}

func Close() {
	for _, c := range cronMap {
		c.Stop()
	}
}
