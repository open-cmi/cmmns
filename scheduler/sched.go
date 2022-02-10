package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/cmmns/storage/rdb"
)

var Sched *Scheduler

// scheduler.hash.namespace.group
// scheduler.zset.namespace.group
// scheduler.zset.namespace.group

// Scheduler scheduler
type Scheduler struct {
	Namespace string
	Mutex     sync.Mutex
	Providers map[string]*Provider
	Consumers map[string]*Consumer
}

func NewScheduler(namespace string) *Scheduler {
	var sched Scheduler

	sched.Namespace = namespace
	sched.Consumers = make(map[string]*Consumer, 0)
	sched.Providers = make(map[string]*Provider, 0)
	return &sched
}

func (s *Scheduler) AddJob(job *Job, option *ProviderOption) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// 添加任务
	cache := rdb.GetCache(rdb.TaskCache)

	var count int = 1
	if option.Type == "anyone" {
		zkey := fmt.Sprintf("scheduler.zset.%s.%s", s.Namespace, option.Group)

		cache.ZAddNX(context.TODO(), zkey, &redis.Z{
			Score:  float64(time.Now().Unix())/1000000 + float64(job.Priority)*10000,
			Member: job.ID,
		})
	} else if option.Type == "everyone" {
		for key, _ := range s.Consumers {
			zkey := fmt.Sprintf("scheduler.zset.%s.%s.%s", s.Namespace, option.Group, key)

			cache.ZAddNX(context.TODO(), zkey, &redis.Z{
				Score:  float64(time.Now().Unix())/1000000 + float64(job.Priority)*10000,
				Member: job.ID,
			})
		}
		count = len(s.Consumers)
	}
	job.State = "Pending"
	job.Count = count

	key := fmt.Sprintf("scheduler.hash.%s.%s.%s", s.Namespace, option.Group, job.ID)
	cache.HSet(context.TODO(), key, "id", job.ID, "type", job.Type, "priority", job.Priority,
		"state", job.State, "count", job.Count, "content", job.Content)

	return nil
}

func (s *Scheduler) GetJob(option *ConsumerOption) *Job {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// 添加任务
	cache := rdb.GetCache(rdb.TaskCache)

	key := fmt.Sprintf("scheduler.zset.%s.%s.%s", s.Namespace, option.Group, option.Identity)
	z, _ := cache.ZPopMax(context.TODO(), key, 1).Result()
	if len(z) == 0 {
		key = fmt.Sprintf("scheduler.zset.%s.%s", s.Namespace, option.Group)
		z, _ = cache.ZPopMax(context.TODO(), key, 1).Result()
		if len(z) == 0 {
			return nil
		}
	}

	jobID := z[0].Member.(string)
	key = fmt.Sprintf("scheduler.hash.%s.%s.%s", s.Namespace, option.Group, jobID)
	jobMap, err := cache.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil
	}
	var job Job
	job.ID, _ = jobMap["id"]
	job.Type, _ = jobMap["type"]
	job.Priority, _ = strconv.Atoi(jobMap["priority"])
	job.State = "Running"
	job.Count, _ = strconv.Atoi(jobMap["count"])
	job.Content, _ = jobMap["content"]

	// 改变job 状态
	key = fmt.Sprintf("scheduler.hash.%s.%s.%s", s.Namespace, option.Group, jobID)
	_, err = cache.HSet(context.TODO(), key, "state", job.State).Result()

	return &job
}

func (s *Scheduler) HasJob(option *ConsumerOption) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// 添加任务
	cache := rdb.GetCache(rdb.TaskCache)

	key := fmt.Sprintf("scheduler.zset.%s.%s", s.Namespace, option.Group)
	count, _ := cache.ZCard(context.TODO(), key).Result()
	if count == 0 {
		key = fmt.Sprintf("scheduler.zset.%s.%s.%s", s.Namespace, option.Group, option.Identity)
		count, _ = cache.ZCard(context.TODO(), key).Result()
		if count == 0 {
			return false
		}
	}
	return true
}

func (sched *Scheduler) NewProvider(option *ProviderOption) *Provider {
	provider := &Provider{
		Sched:  sched,
		Option: option,
	}
	sched.Providers[option.Identity] = provider
	return provider
}

func (sched *Scheduler) GetProvider(identity string) *Provider {
	sched.Mutex.Lock()
	defer sched.Mutex.Unlock()

	provider, ok := sched.Providers[identity]
	if ok {
		return provider
	}
	return nil
}

func (sched *Scheduler) NewConsumer(option *ConsumerOption) *Consumer {

	sched.Mutex.Lock()
	defer sched.Mutex.Unlock()

	consumer, found := sched.Consumers[option.Identity]
	if found {
		return nil
	}

	consumer = &Consumer{
		Sched:  sched,
		Option: option,
	}
	sched.Consumers[option.Identity] = consumer
	return consumer
}

func (sched *Scheduler) GetConsumer(identity string) *Consumer {
	sched.Mutex.Lock()
	defer sched.Mutex.Unlock()

	consumer, ok := sched.Consumers[identity]
	if ok {
		return consumer
	}
	return nil
}

func GetScheduler() *Scheduler {
	return Sched
}

func Init() {
	Sched = NewScheduler("default")
}