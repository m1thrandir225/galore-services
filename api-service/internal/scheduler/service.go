// Package scheduler provides a background task scheduler
package scheduler

// Service Interface for a scheduler service used for background jobs
type Service interface {
	EnqueueJob(jobName string, jobArguments map[string]interface{}) error
	RegisterJob(jobName string, noRetry bool, handler func(args map[string]interface{}) error)
	RegisterCronJob(name, cronSpec string)
	Start()
	Stop()
}
