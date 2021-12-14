package scheduler

import "encoding/json"

// Job job
type Job struct {
	Name      string
	Priority  int
	State     string // init, waiting, appendding, running, completed
	content   JobContent
	scheduler *JobScheduler
}

// Job job struct
type JobContent struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Content interface{} `json:"content,omitempty"`
}

// JobResult job exec result
type JobResult struct {
	JobID   string          `json:"jobid"`
	JobType string          `json:"jobtype"`
	Results json.RawMessage `json:"results"`
}

// NewJobScheduler new job request
func NewJob(priority int, jc JobContent) *Job {
	var job Job

	job.Name = jc.ID
	job.Priority = priority
	job.State = "init"
	job.content = jc
	job.scheduler = NewJobScheduler(jc.ID, priority)
	return &job
}

func (j *Job) SetOption(option string, value interface{}) error {
	return j.scheduler.SetSchedulerOption(option, value)
}

func (j *Job) Schedule(mode string) {
	err := j.scheduler.Schedule(j.content, mode)
	if err == nil {
		j.State = "appendding"
	}
}

func (j *Job) SetNamespace(namespace string) {
	j.scheduler.Namespace = namespace
}

func (j *Job) SetState(state string) {
	j.State = state
}
