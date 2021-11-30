package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/open-cmi/cmmns/db"
)

// Executor job executor
type Executor struct {
	Name     string `json:"name"`
	DeviceID string `json:"deviceid"`
	Address  string `json:"address"`
	Group    int    `json:"group"`
}

func formatExecutorKey(deviceid string) string {
	key := fmt.Sprintf("agent_executor_%s", deviceid)
	return key
}

// RegisterExecutor 注册新的执行器
// name表示执行器的名称
func RegisterExecutor(name string, deviceid string, address string, group int) error {
	var exer Executor

	exer.Name = name
	exer.DeviceID = deviceid
	exer.Address = address
	exer.Group = group

	exeStr, err := json.Marshal(exer)
	if err != nil {
		return err
	}

	cache := db.GetCache(db.AgentCache)
	key := formatExecutorKey(deviceid)
	_, err = cache.Set(context.TODO(), key, exeStr, time.Second*5*60).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetExecutor executor is exist
func GetExecutor(deviceid string) (executor Executor, err error) {
	// 先查缓存是否存在
	key := formatExecutorKey(deviceid)
	cache := db.GetCache(db.AgentCache)
	executorStr, err := cache.Get(context.TODO(), key).Result()
	if err != nil {
		return executor, err
	}
	err = json.Unmarshal([]byte(executorStr), &executor)
	if err != nil {
		return executor, err
	}
	return executor, nil
}

// Unregister 注销执行器
func (e *Executor) Unregister(name string) {
	return
}

// Refresh executor refresh
func (e *Executor) Refresh() error {
	cache := db.GetCache(db.AgentCache)
	key := formatExecutorKey(e.DeviceID)
	_, err := cache.Expire(context.TODO(), key, time.Second*5*60).Result()
	return err
}

// GetAllExecutors get all executors
func GetAllExecutors() (executors []Executor, err error) {
	pattern := "agent_executor_*"
	cache := db.GetCache(db.AgentCache)
	var cursor uint64 = 0
	for {
		keys, cursor, err := cache.Scan(context.TODO(), cursor, pattern, 10).Result()
		if err != nil {
			return []Executor{}, err
		}
		for _, key := range keys {
			executorStr, err := cache.Get(context.TODO(), key).Result()
			fmt.Println(executorStr, err)
			if err != nil {
				continue
			}
			var executor Executor
			err = json.Unmarshal([]byte(executorStr), &executor)
			if err != nil {
				fmt.Println("unmashal", err)
				return []Executor{}, err
			}
			executors = append(executors, executor)
		}

		if cursor == 0 {
			break
		}
	}

	return executors, nil
}
