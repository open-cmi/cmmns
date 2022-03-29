package scheduler

import "github.com/open-cmi/cmmns/common/def"

type ConsumerOption struct {
	Identity string
	Group    string
}

type Consumer struct {
	Sched  *Scheduler
	Option *ConsumerOption
}

func (c *Consumer) GetJob() *def.JobRequest {

	job := c.Sched.jobDequeue(c.Option)
	if job != nil {
		return &def.JobRequest{
			ID:      job.ID,
			Type:    job.Type,
			Content: job.Content,
		}
	}
	return nil
}

func (c *Consumer) HasJob() bool {
	return c.Sched.hasJob(c.Option)
}

func (c *Consumer) JobDone(jobResp *def.JobResponse) {
	job := Get(nil, "id", jobResp.ID)
	if job == nil {
		return
	}

	job.Code = jobResp.Code
	job.Msg = jobResp.Msg
	job.Result = jobResp.Result
	c.Sched.jobDone(c.Option, job)
}
