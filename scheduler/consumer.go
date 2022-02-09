package scheduler

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
