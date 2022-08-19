package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/open-cmi/cmmns/essential/rdb"
)

func GetCache(devID string) *Model {
	cache := rdb.GetClient("public")
	if cache == nil {
		return nil
	}

	key := fmt.Sprintf("agent-%s", devID)
	v, err := cache.Get(context.TODO(), key).Result()
	if err != nil {
		return nil
	}
	var mdl Model
	err = json.Unmarshal([]byte(v), &mdl)
	if err != nil {
		return nil
	}
	go cache.Expire(context.TODO(), key, 1*time.Minute).Result()

	return &mdl
}

func SetCache(mdl *Model) error {
	cache := rdb.GetClient("public")
	if cache == nil {
		return errors.New("cache not exist")
	}

	v, err := json.Marshal(mdl)
	if err != nil {
		return nil
	}
	_, err = cache.Set(context.TODO(),
		fmt.Sprintf("agent-%s", mdl.DevID),
		string(v), 1*time.Minute).Result()
	return err
}

func RefreshCache(devID string) error {
	cache := rdb.GetClient("public")
	if cache == nil {
		return errors.New("cache not exist")
	}
	_, err := cache.Expire(context.TODO(),
		fmt.Sprintf("agent-%s", devID),
		1*time.Minute).Result()
	return err
}

func DeleteCache(devID string) error {
	cache := rdb.GetClient("public")
	if cache == nil {
		return errors.New("cache not exist")
	}
	_, err := cache.Del(context.TODO(),
		fmt.Sprintf("agent-%s", devID),
	).Result()
	return err
}
