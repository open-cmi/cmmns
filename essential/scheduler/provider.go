package scheduler

import (
	"errors"
	"time"
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

func (p *Provider) AddJob(jobReq *JobRequest) error {
	// 这里做job的校验
	if jobReq.ID == "" || jobReq.Type == "" {
		return errors.New("job id and type should not be empty")
	}
	job := new(Job)
	job.StartedTime = time.Now().Unix()
	job.IsNew = true
	job.ID = jobReq.ID
	job.Type = jobReq.Type
	job.Priority = jobReq.Priority
	job.State = "Pending"
	job.RunType = jobReq.RunType
	job.RunSpec = jobReq.RunSpec
	job.SchedGroup = p.Option.Group
	job.SchedObject = p.Option.Type
	job.Content = jobReq.Content

	if job.SchedObject == "everyone" {
		for _, cons := range p.Sched.Consumers {
			if cons.Option.Group == job.SchedGroup {
				job.Count++
			}
		}
	} else {
		job.Count = 1
	}
	err := p.Sched.addJob(job)
	return err
}
