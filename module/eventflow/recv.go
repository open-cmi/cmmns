package eventflow

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/rdb"
)

var gEventCallbacks map[string]func(v string) = make(map[string]func(v string))

func EventFlowRecvProc(key string, channel string) {
	logger.Debugf("EventFlowRecvProc, key %s channel %s\n", key, channel)
	callback, ok := gEventCallbacks[key]
	if !ok {
		// 如果未注册的话，就直接退出
		return
	}

	if GetConf().Relay == "rdb" {
		rcli := rdb.GetClient(0)
		if GetConf().Method == "queue" {
			for {
				v, err := rcli.Conn.LPop(context.TODO(), channel).Result()
				if err != nil {
					if err != redis.Nil {
						// 错误，重新连接
						rcli.Reconnect()
					}
					// 数据为空，睡眠
					time.Sleep(1 * time.Second)
					continue
				}
				callback(v)
			}
		}
	}
}

func RegisterChannel(key string, f func(v string)) error {
	_, ok := gEventCallbacks[key]
	if ok {
		return fmt.Errorf("event stream %s is existing", key)
	}

	gEventCallbacks[key] = f
	return nil
}
