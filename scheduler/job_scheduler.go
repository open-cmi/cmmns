package scheduler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/open-cmi/cmmns/storage/rdb"

	"github.com/robfig/cron/v3"
)

const (
	// SchedulerCommon common job，通用任务，任何一个执行器都能执行
	SchedulerCommon = 1

	// SchedulerAppointed to self job list，委任给某些执行器执行
	SchedulerAppointed = 2

	// SchedulerGroup to group，仅某个组的执行器才能执行的任务
	SchedulerGroup = 3
)

// JobScheduler job scheduler
type JobScheduler struct {
	Namespace string     `json:"namespace"` // scheduler name
	Name      string     `json:"name"`      // 任务name，同一个scheduler下，不同的jobscheduler，name不能相同
	Priority  int        `json:"priority"`  // 任务优先级，值越小，优先级越高
	Object    int        `json:"object"`    // 调度对象，是给指定的执行器执行，还是给group运行,1表示给指定的执行器，2表示给group
	Executors []Executor `json:"executors"` // 如果是给某个执行器执行，必须指定该项
	Group     int        `json:"group"`     // 指定运行组
	RunMode   string     `json:"runmode"`   // runmode ,一次运行，周期运行，定时运行
	Spec      string     `json:"spec"`      // 定时器格式
}

// NewJobScheduler new job request
func NewJobScheduler(jobname string, priority int) *JobScheduler {
	var jobSched JobScheduler

	jobSched.Name = jobname
	jobSched.Priority = priority
	jobSched.Namespace = "global"
	return &jobSched
}

// SetSchedulerOption 设置选项
func (s *JobScheduler) SetSchedulerOption(option string, value interface{}) error {
	if option == "GroupJob" {
		group, ok := value.(int)
		if !ok {
			return errors.New("group not valid")
		}
		s.Object = SchedulerGroup
		s.Group = group
	} else if option == "CommonJob" {
		s.Object = SchedulerCommon
	} else if option == "AppointedJob" {
		s.Object = SchedulerAppointed
		executors, ok := value.([]Executor)
		if !ok {
			return errors.New("executors not valid")
		}
		s.Executors = executors
	}

	return nil
}

// Schedule 加入调度
func (s *JobScheduler) Schedule(req JobContent, runMode string) error {
	if runMode == "once" {
		if s.Object == SchedulerAppointed {
			for _, executor := range s.Executors {
				err := s.EnqueueAppointed(&executor, req)
				if err != nil {
					return err
				}
			}
		} else if s.Object == SchedulerGroup {
			err := s.EnqueueGroup(s.Group, req)
			if err != nil {
				return err
			}
		}
	} else {
		c := cron.New(cron.WithSeconds())

		c.AddFunc(s.Spec, func() {
			if s.Object == SchedulerAppointed {
				for _, executor := range s.Executors {
					err := s.EnqueueAppointed(&executor, req)
					if err != nil {
						return
					}
				}
			} else if s.Object == SchedulerGroup {
				err := s.EnqueueGroup(s.Group, req)
				if err != nil {
					return
				}
			}
		})
		c.Start()
	}
	return nil
}

// EnqueueAppointed push task to single
func (r *JobScheduler) EnqueueAppointed(executor *Executor, job JobContent) error {
	cache := rdb.GetCache(rdb.TaskCache)

	jsonTask, err := json.Marshal(job)
	if err != nil {
		return err
	}

	queue := fmt.Sprintf("%s_appointed_job_%s", r.Namespace, executor.DeviceID)
	err = cache.LPush(context.TODO(), queue, string(jsonTask)).Err()
	if err != nil {
		return err
	}

	return nil
}

// EnqueueGroup push task to group
func (r *JobScheduler) EnqueueGroup(group int, job JobContent) error {
	cache := rdb.GetCache(rdb.TaskCache)

	jsonTask, err := json.Marshal(job)
	if err != nil {
		return err
	}
	queue := fmt.Sprintf("%s_group_job_%d", r.Namespace, group)
	err = cache.LPush(context.TODO(), queue, string(jsonTask)).Err()
	if err != nil {
		return err
	}

	return nil
}
