package scheduler

import (
	"errors"
	"time"
)

var ShouldStop bool = false
var SchedMap map[string]*Scheduler

// GetJob get task
func GetJob(executor *Executor) (job JobContent, err error) {
	for _, sched := range SchedMap {
		job, err = sched.GetJob(executor)
		if err == nil {
			return job, nil
		}
	}
	return job, errors.New("no jobs")
}

// HasJob check whether executor has job
func HasJob(executor *Executor) bool {
	for _, sched := range SchedMap {
		if sched.HasJob(executor) {
			return true
		}
	}
	return false
}

func DeInit() {
	ShouldStop = true
}

func Init() {
	go func() {
		for !ShouldStop {
			for _, sched := range SchedMap {
				sched.Run()
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func init() {
	SchedMap = make(map[string]*Scheduler, 0)
}
