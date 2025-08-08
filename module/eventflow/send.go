package eventflow

import (
	"context"

	"github.com/open-cmi/gobase/essential/logger"
)

func EventFlowSend(chankey string, v string) {
	channel, ok := GetConf().Channels[chankey]
	if !ok {
		logger.Errorf("chankey %s not found\n", chankey)
		return
	}
	if GetConf().Relay == "rdb" {
		switch GetConf().Method {
		case "push":
			logger.Debugf("push %s: %s\n", channel, v)
			_, err := rcli.Conn.RPush(context.TODO(), channel, v).Result()
			if err != nil {
				logger.Errorf("push %s to redis failed: %s\n", channel, err.Error())
				rcli.Reconnect()
			}
		case "publish":
			logger.Debugf("publish %s: %s\n", channel, v)
			_, err := rcli.Conn.Publish(context.TODO(), channel, v).Result()
			if err != nil {
				logger.Errorf("publish %s to redis failed: %s\n", channel, err.Error())
				rcli.Reconnect()
			}
		}
	}
}

func EventFlowMultiSend(chankey string, datas []interface{}) {
	channel, ok := GetConf().Channels[chankey]
	if !ok {
		logger.Errorf("chankey %s not found\n", chankey)
		return
	}

	if GetConf().Relay == "rdb" {
		switch GetConf().Method {
		case "queue":
			logger.Debugf("batch push %s\n", channel)
			_, err := rcli.Conn.RPush(context.TODO(), channel, datas...).Result()
			if err != nil {
				logger.Errorf("push %s to redis failed: %s\n", channel, err.Error())
				rcli.Reconnect()
			}
		case "pubsub":
			logger.Debugf("batch publish %s\n", channel)
			for _, v := range datas {
				_, err := rcli.Conn.Publish(context.TODO(), channel, v).Result()
				if err != nil {
					logger.Errorf("publish %s to redis failed: %s\n", channel, err.Error())
					rcli.Reconnect()
				}
			}
		}
	}
}
