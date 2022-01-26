package notification

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/cmmns/storage/rdb"
)

// parseMsg
func dispatchMsg(msgchannel string, msg string) {
	switch msgchannel {
	case "Deploy":
	default:
	}
}

// Init function
func Init() {
	cache := rdb.GetCache(rdb.TaskCache)
	pubsub := cache.Subscribe(context.TODO(),
		"Deploy",
	)

	for {
		iface, _ := pubsub.Receive(context.TODO())

		switch msg := iface.(type) {
		case *redis.Subscription:
			fmt.Println("recv Subscription")
		case *redis.Pong:
			fmt.Println("recv Pong")
		case *redis.Message:
			fmt.Printf("recv Message %s, %s\n", msg.Channel, msg.Payload)
			dispatchMsg(msg.Channel, msg.Payload)
		}
	}
}
