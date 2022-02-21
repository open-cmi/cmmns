package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/cmmns/essential/logger"
	"github.com/open-cmi/cmmns/essential/storage/rdb"
)

var Sched *Scheduler

// scheduler.hash.namespace.group
// scheduler.zset.namespace.group
// scheduler.zset.namespace.group

type JobDoneFunc func(*Job, *ConsumerOption)
type JobProgressFunc func(*Job, *ConsumerOption)
type JobTimeoutFunc func(*Job, *ConsumerOption)

type JobCallback struct {
	Done     JobDoneFunc
	Progress JobProgressFunc
	Timeout  JobTimeoutFunc
}

// Scheduler scheduler
type Scheduler struct {
	Namespace string
	Mutex     sync.Mutex
	Providers map[string]*Provider
	Consumers map[string]*Consumer
	Callback  map[string]JobCallback
}

func NewScheduler(namespace string) *Scheduler {
	var sched Scheduler

	sched.Namespace = namespace
	sched.Consumers = make(map[string]*Consumer, 0)
	sched.Providers = make(map[string]*Provider, 0)
	sched.Callback = make(map[string]JobCallback)
	return &sched
}

func (s *Scheduler) enqueue(job *Job) error {
	cache := rdb.GetCache("scheduler")

	var err error
	if job.SchedObject == "anyone" {
		zkey := fmt.Sprintf("scheduler.zset.%s.%s", s.Namespace, job.SchedGroup)

		_, err = cache.ZAddNX(context.TODO(), zkey, &redis.Z{
			Score:  float64(time.Now().Unix())/1000000 + float64(job.Priority)*10000,
			Member: job.ID,
		}).Result()
	} else if job.SchedObject == "everyone" {
		for key, cons := range s.Consumers {
			if cons.Option.Group == job.SchedGroup {
				zkey := fmt.Sprintf("scheduler.zset.%s.%s.%s", s.Namespace, job.SchedGroup, key)

				_, err = cache.ZAddNX(context.TODO(), zkey, &redis.Z{
					Score:  float64(time.Now().Unix())/1000000 + float64(job.Priority)*10000,
					Member: job.ID,
				}).Result()

				if err != nil {
					logger.Errorf("zaddnx failed: %s\n", err.Error())
					break
				}
			}
		}
	}

	job.State = "Pending"
	return err
}

func (s *Scheduler) dequeue(option *ConsumerOption) (id string) {
	// 获取cache
	cache := rdb.GetCache("scheduler")
	if cache == nil {
		return ""
	}

	// 从优先级集合中获取jobid
	key := fmt.Sprintf("scheduler.zset.%s.%s.%s", s.Namespace, option.Group, option.Identity)
	z, _ := cache.ZPopMax(context.TODO(), key, 1).Result()
	if len(z) == 0 {
		key = fmt.Sprintf("scheduler.zset.%s.%s", s.Namespace, option.Group)
		z, _ = cache.ZPopMax(context.TODO(), key, 1).Result()
		if len(z) == 0 {
			return ""
		}
	}

	jobID := z[0].Member.(string)
	return jobID
}

func (s *Scheduler) addJob(job *Job) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	err := job.Save()
	if err == nil {
		// 如果保存失败，需要回退?
		err = s.enqueue(job)
	}

	return err
}

func (s *Scheduler) getJob(option *ConsumerOption) *JobRequest {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// 获取cache
	cache := rdb.GetCache("scheduler")
	if cache == nil {
		return nil
	}
	jobID := s.dequeue(option)

	if jobID != "" {
		var jobRequest JobRequest
		job := Get("id", jobID)
		if job == nil {
			return nil
		}
		jobRequest.ID = job.ID
		jobRequest.Type = job.Type
		jobRequest.Content = job.Content

		job.State = Running

		job.Save()

		return &jobRequest
	}
	return nil
}

func (s *Scheduler) RegisterJobCallback(jobType string, callback JobCallback) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, found := s.Callback[jobType]
	if found {
		return errors.New("job callback has been registered")
	}
	s.Callback[jobType] = callback
	return nil
}

func (s *Scheduler) jobDone(option *ConsumerOption, jobResp *JobResponse) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// 添加cache
	cache := rdb.GetCache("scheduler")
	if cache == nil {
		return
	}

	job := Get("id", jobResp.ID)
	if job == nil {
		return
	}

	job.Result = jobResp.Result
	job.Done += 1
	if job.Count == job.Done {
		job.State = "Done"
		job.StoppedTime = time.Now().Unix()
		job.Save()
		callback, found := s.Callback[jobResp.Type]
		if found {
			callback.Done(job, option)
		}
	} else {
		// 只增加done的数量，不删除数据
		job.Save()
		callback, found := s.Callback[jobResp.Type]
		if found {
			callback.Progress(job, option)
		}
	}
}

func (s *Scheduler) hasJob(option *ConsumerOption) bool {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// 添加任务
	cache := rdb.GetCache("scheduler")

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
	_, found := sched.Providers[option.Identity]
	if found {
		return nil
	}
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

	_, found := sched.Consumers[option.Identity]
	if found {
		return nil
	}

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

	consumer, ok := sched.Consumers[identity]
	if ok {
		return consumer
	}
	return nil
}

func GetScheduler() *Scheduler {
	if Sched == nil {
		Sched = NewScheduler("default")
	}
	return Sched
}
