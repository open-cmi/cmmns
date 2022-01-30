package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/cmmns/storage/rdb"
)

var Sched *Scheduler

// scheduler.hash.namespace.group
// scheduler.zset.namespace.group
// scheduler.zset.namespace.group.consumer

// Scheduler scheduler
type Scheduler struct {
	Namespace string
	Mutex     sync.Mutex
	Providers map[string]*Provider
	Consumers map[string]*Consumer
}

type ProviderOption struct {
	Identity   string
	Type       string
	Group      string
	ConsumerID string
}

type ConsumerOption struct {
	Identity string
	Group    string
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

	// 添加任务
	cache := rdb.GetCache(rdb.TaskCache)

	if option.Type == "anyone" {
		zkey := fmt.Sprintf("scheduler.zset.%s.%s", s.Namespace, option.Group)

		cache.ZAddNX(context.TODO(), zkey, &redis.Z{
			Score:  float64(time.Now().Unix()) / 1000,
			Member: job.ID,
		})
	}
	job.State = "Pending"
	jobJson, err := json.Marshal(*job)
	if err != nil {
		s.Mutex.Unlock()
		return err
	}

	key := fmt.Sprintf("scheduler.hash.%s.%s", s.Namespace, option.Group)
	cache.HSet(context.TODO(), key, job.ID, jobJson)
	s.Mutex.Unlock()
	return nil
}

func (s *Scheduler) GetJob(option *ConsumerOption) *Job {
	s.Mutex.Lock()

	// 添加任务
	cache := rdb.GetCache(rdb.TaskCache)

	key := fmt.Sprintf("scheduler.zset.%s.%s", s.Namespace, option.Group)
	z, _ := cache.ZPopMax(context.TODO(), key, 1).Result()
	if len(z) == 0 {
		s.Mutex.Unlock()
		return nil
	}
	jobID := z[0].Member.(string)
	key = fmt.Sprintf("scheduler.hash.%s.%s", s.Namespace, option.Group)
	jsonTask, err := cache.HGet(context.TODO(), key, jobID).Result()
	if err != nil {
		s.Mutex.Unlock()
		return nil
	}
	var job Job
	err = json.Unmarshal([]byte(jsonTask), &job)
	if err != nil {
		s.Mutex.Unlock()
		return nil
	}
	s.Mutex.Unlock()
	return &job
}

func (s *Scheduler) HasJob(option *ConsumerOption) bool {
	s.Mutex.Lock()

	// 添加任务
	cache := rdb.GetCache(rdb.TaskCache)

	key := fmt.Sprintf("scheduler.zset.%s.%s", s.Namespace, option.Group)
	z, _ := cache.ZPopMax(context.TODO(), key, 1).Result()
	if len(z) == 0 {
		s.Mutex.Unlock()
		return false
	}
	s.Mutex.Unlock()
	return true
}

type Provider struct {
	Sched  *Scheduler
	Option *ProviderOption
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

	return sched.Providers[identity]
}

func (p *Provider) AddJob(job *Job) error {

	err := p.Sched.AddJob(job, p.Option)
	return err
}

type Consumer struct {
	Sched  *Scheduler
	Option *ConsumerOption
}

func (sched *Scheduler) NewConsumer(option *ConsumerOption) *Consumer {

	consumer := &Consumer{
		Sched:  sched,
		Option: option,
	}
	sched.Consumers[option.Identity] = consumer
	return consumer
}

func (sched *Scheduler) GetConsumer(identity string) *Consumer {
	sched.Mutex.Lock()
	defer sched.Mutex.Unlock()

	return sched.Consumers[identity]
}

func (c *Consumer) GetJob() *Job {

	job := c.Sched.GetJob(c.Option)
	return job
}

func (c *Consumer) HasJob() bool {
	return c.Sched.HasJob(c.Option)
}

// Job job
type Job struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"`
	Priority int         `json:"priority"`
	State    string      `json:"state"`
	Content  interface{} `json:"content,omitempty"`
}

func GetScheduler() *Scheduler {
	return Sched
}

func Init() {
	Sched = NewScheduler("global")
}
