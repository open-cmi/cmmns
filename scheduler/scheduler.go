package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-cmi/cmmns/storage/rdb"
)

type Scheduler struct {
	ShouldStop bool
	Namespace  string
	Jobs       map[string]*Job
}

func NewScheduler(namespace string) *Scheduler {
	var sched Scheduler

	sched.Namespace = namespace
	sched.ShouldStop = false
	sched.Jobs = make(map[string]*Job, 0)
	return &sched
}

func (s *Scheduler) AddJob(job *Job) {
	job.SetNamespace(s.Namespace)
	job.SetState("waiting")
	s.Jobs[job.Name] = job
}

func (s *Scheduler) Setup() {
	SchedMap[s.Namespace] = s
}

func (s *Scheduler) Run() {
	// 循环执行Scheduler，将添加的任务，放入到队列中，并取回结果
	var count int = 0
	for _, job := range s.Jobs {
		fmt.Println("job run: ", job)
		if job.State == "waiting" {
			count++
			job.Schedule("once")
		} else if job.State == "completed" {
			// 任务完成后，如何处理
		}
	}
}

// CheckAppointedJob check self job
func (r *Scheduler) HasAppointedJob(executor *Executor) bool {
	queue := fmt.Sprintf("%s_appointed_job_%s", r.Namespace, executor.DeviceID)
	cache := rdb.GetCache(rdb.TaskCache)
	l, err := cache.LLen(context.TODO(), queue).Result()
	if err != nil || l == 0 {
		return false
	}

	return true
}

// CheckGroupJob check self job
func (r *Scheduler) HasGroupJob(executor *Executor) bool {
	queue := fmt.Sprintf("%s_group_job_%d", r.Namespace, executor.Group)
	cache := rdb.GetCache(rdb.TaskCache)
	l, err := cache.LLen(context.TODO(), queue).Result()
	if err != nil || l == 0 {
		return false
	}

	return true
}

func (s *Scheduler) HasJob(executor *Executor) bool {
	if s.HasAppointedJob(executor) {
		return true
	}
	if s.HasGroupJob(executor) {
		return true
	}
	return false
}

// GetGroupJob pop group job
func (r *Scheduler) GetGroupJob(executor *Executor) (req JobContent, err error) {
	queue := fmt.Sprintf("%s_group_job_%d", r.Namespace, executor.Group)
	cache := rdb.GetCache(rdb.TaskCache)
	taskstr, err := cache.LPop(context.TODO(), queue).Result()
	if err != nil {
		return req, err
	}
	err = json.Unmarshal([]byte(taskstr), &req)
	if err != nil {
		return req, err
	}
	return req, nil
}

// GetAppointedJob pop self task
func (r *Scheduler) GetAppointedJob(executor *Executor) (req JobContent, err error) {
	queue := fmt.Sprintf("%s_appointed_job_%s", r.Namespace, executor.DeviceID)
	cache := rdb.GetCache(rdb.TaskCache)
	taskstr, err := cache.LPop(context.TODO(), queue).Result()
	if err != nil {
		return req, err
	}
	err = json.Unmarshal([]byte(taskstr), &req)
	if err != nil {
		return req, err
	}
	return req, nil
}

// GetJob pop task
func (r *Scheduler) GetJob(executor *Executor) (JobContent, error) {
	taskstr, err := r.GetAppointedJob(executor)
	if err != nil {
		taskstr, err = r.GetGroupJob(executor)
	}

	return taskstr, err
}
