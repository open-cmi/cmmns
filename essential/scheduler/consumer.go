package scheduler

import "github.com/open-cmi/cmmns/common/job"

type ConsumerOption struct {
	Identity string
	Group    string
}

type Consumer struct {
	Sched  *Scheduler
	Option *ConsumerOption
}

func (c *Consumer) GetJob() *Job {

	job := c.Sched.GetJob(c.Option)
	return job
}

func (c *Consumer) HasJob() bool {
	return c.Sched.HasJob(c.Option)
}

func (c *Consumer) JobDone(job *job.Response) {
	c.Sched.JobDone(c.Option, job)
}
