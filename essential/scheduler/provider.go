package scheduler

type ProviderOption struct {
	Identity string
	Type     string
	Group    string
}

type Provider struct {
	Sched  *Scheduler
	Option *ProviderOption
}

func (p *Provider) AddJob(job *Job) error {

	err := p.Sched.AddJob(job, p.Option)
	return err
}
