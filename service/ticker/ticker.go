package ticker

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
)

type TickerFunc func()

type Ticker struct {
	Spec string
	Func TickerFunc
}

var tickers map[string]Ticker = make(map[string]Ticker)

func Register(name string, spec string, f TickerFunc) error {
	_, found := tickers[name]
	if found {
		errMsg := fmt.Sprintf("ticker %s registered failed", name)
		return errors.New(errMsg)
	}
	tickers[name] = Ticker{
		Spec: spec,
		Func: f,
	}
	return nil
}

var cronInstance *cron.Cron

func Init() error {
	cronInstance = cron.New(cron.WithSeconds())

	for _, ticker := range tickers {
		cronInstance.AddFunc(ticker.Spec, ticker.Func)
	}

	return nil
}

func Run() error {
	cronInstance.Start()
	return nil
}

func Close() {
	cronInstance.Stop()
}
