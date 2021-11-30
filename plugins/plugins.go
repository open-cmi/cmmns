package plugins

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/cmmns/db"
	"github.com/open-cmi/cmmns/plugins/deploy"
)

// parseMsg
func dispatchMsg(msgchannel string, msg string) {
	switch msgchannel {
	case "DeployAgent":
		go deploy.Exec(msg)
	default:
	}
}

// Init function
func Init() {
	cache := db.GetCache(db.TaskCache)
	pubsub := cache.Subscribe(context.TODO(),
		"DeployAgent",
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
