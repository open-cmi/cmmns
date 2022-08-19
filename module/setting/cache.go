package setting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/open-cmi/cmmns/essential/rdb"
)

func GetCache(key string) *Model {
	cache := rdb.GetClient("public")
	if cache == nil {
		return nil
	}

	rkey := fmt.Sprintf("setting.%s", key)
	v, err := cache.Get(context.TODO(), rkey).Result()
	if err != nil {
		return nil
	}
	var mdl Model
	err = json.Unmarshal([]byte(v), &mdl)
	if err != nil {
		return nil
	}

	go cache.Expire(context.TODO(), key,
		30*time.Minute).Result()

	return &mdl
}

func SetCache(key string, mdl *Model) error {
	cache := rdb.GetClient("public")
	if cache == nil {
		return errors.New("cache not exist")
	}

	v, err := json.Marshal(mdl)
	if err != nil {
		return nil
	}
	_, err = cache.Set(context.TODO(),
		fmt.Sprintf("setting-%s", key),
		string(v), 1*time.Minute).Result()
	return err
}

func DelCache(key string) error {
	cache := rdb.GetClient("public")
	if cache == nil {
		return errors.New("cache not exist")
	}

	_, err := cache.Del(context.TODO(),
		fmt.Sprintf("setting-%s", key)).Result()
	return err
}
