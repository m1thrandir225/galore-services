package background_jobs

type SchedulerService interface {
	EnqueueJob(jobName string, jobArguments map[string]interface{}) error
	RegisterJob(jobName string, handler func(args map[string]interface{}) error)
	Start()
	Stop()
}
