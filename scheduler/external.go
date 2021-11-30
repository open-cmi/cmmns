package scheduler

// GetTask get task
func GetTask(executor Executor) (Request, error) {
	return Sched.GetTask(executor)
}

// CheckAvailableJob check available job
func CheckAvailableJob(executor Executor) int {
	return Sched.CheckAvailableJob(executor)
}

// Schedule schedule
func Schedule(jr *JobResource, req Request) {
	Sched.Schedule(jr, req)
}

// Init  init
func init() {
	Sched = NewScheduler()
}
