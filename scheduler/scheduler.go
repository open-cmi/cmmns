package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/open-cmi/cmmns/db"

	"github.com/robfig/cron/v3"
)

// Sched global
var Sched *Scheduler = nil

const (
	// SchedulerNone none job
	SchedulerNone = 0

	// SchedulerSelf to self job list
	SchedulerSelf = 1

	// SchedulerGroup to group
	SchedulerGroup = 2
)

// JobResource job resource
type JobResource struct {
	Name           string     `json:"scheduler_name"`
	ScheduleObject int        `json:"schedule_object"` // 调度对象，是给指定的执行器执行，还是给group运行,1表示给指定的执行器，2表示给group
	Executors      []Executor `json:"executors"`       // 如果是给某个执行器执行，必须指定该项
	Group          int        `json:"group"`           // 指定运行组
	RunMode        string     `json:"runmode"`         // runmode ,一次运行，周期运行，定时运行
	Spec           string     `json:"spec"`            // 定时器格式
}

// NewJobResource new job resource
func NewJobResource(name string, scheduleObject int) *JobResource {
	var jr JobResource
	jr.Name = name
	jr.ScheduleObject = scheduleObject

	return &jr
}

// SetExecutor 添加调度器
func (r *JobResource) SetExecutor(executors []Executor) {
	for _, executor := range executors {
		r.Executors = append(r.Executors, executor)
	}
	return
}

// SetGroup 添加调度器
func (r *JobResource) SetGroup(group int) {
	r.Group = group
	return
}

// Scheduler job scheduler
type Scheduler struct {
	// storage可以在这里设置
	RequestMap  map[string]Request
	ResourceMap map[string]JobResource
	CronMap     map[string]*cron.Cron
}

// NewScheduler new job request
func NewScheduler() *Scheduler {
	var jobr Scheduler
	jobr.RequestMap = make(map[string]Request, 1)
	jobr.ResourceMap = make(map[string]JobResource, 1)
	jobr.CronMap = make(map[string]*cron.Cron, 1)
	return &jobr
}

// Schedule 加入调度
func (r *Scheduler) Schedule(jr *JobResource, req Request) {
	if jr.RunMode == "once" {
		if jr.ScheduleObject == SchedulerSelf {
			for _, executor := range jr.Executors {
				r.EnqueueSelf(executor.DeviceID, req)
			}
		} else if jr.ScheduleObject == SchedulerGroup {
			r.EnqueueGroup(jr.Group, req)
		}
	} else {
		c := cron.New(cron.WithSeconds())

		c.AddFunc(jr.Spec, func() {
			if jr.ScheduleObject == SchedulerSelf {
				for _, executor := range jr.Executors {
					r.EnqueueSelf(executor.DeviceID, req)
				}
			} else if jr.ScheduleObject == SchedulerGroup {
				r.EnqueueGroup(jr.Group, req)
			}
		})
		c.Start()
		r.RequestMap[req.JobID] = req
		r.ResourceMap[req.JobID] = *jr
		r.CronMap[req.JobID] = c
	}

	return
}

// EnqueueSelf push task to single
func (r *Scheduler) EnqueueSelf(deviceid string, job Request) error {
	cache := db.GetCache(db.TaskCache)

	jsonTask, err := json.Marshal(job)
	if err != nil {
		return err
	}
	queue := fmt.Sprintf("self_task_%s", deviceid)
	err = cache.LPush(context.TODO(), queue, string(jsonTask)).Err()
	if err != nil {
		return err
	}

	return nil
}

// EnqueueGroup push task to group
func (r *Scheduler) EnqueueGroup(group int, job Request) error {
	cache := db.GetCache(db.TaskCache)

	jsonTask, err := json.Marshal(job)
	if err != nil {
		return err
	}
	queue := fmt.Sprintf("group_task_%d", group)
	err = cache.LPush(context.TODO(), queue, string(jsonTask)).Err()
	if err != nil {
		return err
	}

	return nil
}

// CheckSchedulerSelf check self task
func (r *Scheduler) CheckSchedulerSelf(executor Executor) bool {
	queue := fmt.Sprintf("self_task_%s", executor.DeviceID)
	cache := db.GetCache(db.TaskCache)
	l, err := cache.LLen(context.TODO(), queue).Result()
	if err != nil || l == 0 {
		return false
	}

	return true
}

// CheckSchedulerGroup check self task
func (r *Scheduler) CheckSchedulerGroup(executor Executor) bool {
	queue := fmt.Sprintf("group_task_%d", executor.Group)
	cache := db.GetCache(db.TaskCache)
	l, err := cache.LLen(context.TODO(), queue).Result()
	if err != nil || l == 0 {
		return false
	}

	return true
}

// CheckAvailableJob check available job
func (r *Scheduler) CheckAvailableJob(executor Executor) int {
	if r.CheckSchedulerSelf(executor) {
		return SchedulerSelf
	}
	if r.CheckSchedulerGroup(executor) {
		return SchedulerGroup
	}

	return SchedulerNone
}

// GetGroupTask pop group task
func (r *Scheduler) GetGroupTask(executor Executor) (req Request, err error) {
	queue := fmt.Sprintf("group_task_%d", executor.Group)
	cache := db.GetCache(db.TaskCache)
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

// GetSelfTask pop self task
func (r *Scheduler) GetSelfTask(executor Executor) (req Request, err error) {
	queue := fmt.Sprintf("self_task_%s", executor.DeviceID)
	cache := db.GetCache(db.TaskCache)
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

// GetTask pop task
func (r *Scheduler) GetTask(executor Executor) (Request, error) {
	taskstr, err := r.GetSelfTask(executor)
	if err != nil {
		taskstr, err = r.GetGroupTask(executor)
	}
	return taskstr, err
}
