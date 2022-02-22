package scheduler

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	RunTypeOnce  = "once"
	RunTypeTimer = "timer"
)

type ProviderOption struct {
	Identity string
	Type     string
	Group    string
}

type Provider struct {
	Sched  *Scheduler
	Option *ProviderOption
}

func (p *Provider) AddJob(jobType string, priority int, content string, runType string, runSpec string) error {
	// 这里做job的校验
	if jobType == "" || runType == "" {
		return errors.New("job type and run type should not be empty")
	}

	job := new(Job)
	job.ID = uuid.NewString()
	job.StartedTime = time.Now().Unix()
	job.IsNew = true
	job.Type = jobType
	job.Priority = priority
	job.State = "Pending"
	job.RunType = runType
	job.RunSpec = runSpec
	job.SchedGroup = p.Option.Group
	job.SchedObject = p.Option.Type
	job.Content = content

	if job.SchedObject == "everyone" {
		for _, cons := range p.Sched.Consumers {
			if cons.Option.Group == job.SchedGroup {
				job.Count++
			}
		}
	} else {
		job.Count = 1
	}
	err := p.Sched.jobEnqueue(job)
	return err
}
