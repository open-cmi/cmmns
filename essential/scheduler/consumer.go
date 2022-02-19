package scheduler

type ConsumerOption struct {
	Identity string
	Group    string
}

type Consumer struct {
	Sched  *Scheduler
	Option *ConsumerOption
}

func (c *Consumer) GetJob() *JobRequest {

	job := c.Sched.getJob(c.Option)
	return job
}

func (c *Consumer) HasJob() bool {
	return c.Sched.hasJob(c.Option)
}

func (c *Consumer) JobDone(job *JobResponse) {
	c.Sched.jobDone(c.Option, job)
}
