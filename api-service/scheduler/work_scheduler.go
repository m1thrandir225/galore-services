package scheduler

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Description:
// GoWork background scheduler implementation of the SchedulerService
//
// Parameters:
// enqueuer: *work.Enqueuer
// workerPool: *work.WorkerPool
type GoworkScheduler struct {
	enqueuer   *work.Enqueuer
	workerPool *work.WorkerPool
}

type JobContext struct{}

// Description:
// Returns a new object of a GoworkScheduler
//
// Parameters:
// namespace: string
// redisUrl: string
//
// Return:
// *GoworkScheduler
func NewGoworkScheduler(namespace string, redisUrl string) *GoworkScheduler {
	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisUrl)
		},
	}
	workerPool := work.NewWorkerPool(JobContext{}, 10, namespace, redisPool)

	enqueuer := work.NewEnqueuer(namespace, redisPool)
	return &GoworkScheduler{
		enqueuer:   enqueuer,
		workerPool: workerPool,
	}
}

// Description
// Enqueue a background job in the Gowork Scheduler
//
// Parameters:
// jobName: string
// jobArguments: map[string]interface{}
//
// Return:
// error
func (p *GoworkScheduler) EnqueueJob(jobName string, jobArguments map[string]interface{}) error {
	_, err := p.enqueuer.Enqueue(jobName, jobArguments)
	return err
}

// Description:
// Register a background job with a CRON
//
// Parameters:
// name: string
// cronSpec: string
func (p *GoworkScheduler) RegisterCronJob(name, cronSpec string) {
	p.workerPool.PeriodicallyEnqueue(cronSpec, name)
}

// Description:
// Register a new job that can be enqueued
//
// Parameters:
// jobName: string
// noRetry: bool
// handler: func(args: map[string]interface) error
func (p *GoworkScheduler) RegisterJob(jobName string, noRetry bool, handler func(args map[string]interface{}) error) {
	var jobOptions work.JobOptions

	if noRetry {
		jobOptions = work.JobOptions{
			MaxFails: 1,
		}
	} else {
		jobOptions = work.JobOptions{
			MaxFails: 0,
		}
	}

	p.workerPool.JobWithOptions(jobName, jobOptions, func(job *work.Job) error {
		args := make(map[string]interface{})
		for key, val := range job.Args {
			args[key] = val
		}
		return handler(args)
	})
}

// Description:
// Start a worker pool from the GoworkScheduler
func (p *GoworkScheduler) Start() {
	log.Println("Starting Go Work scheduler")
	p.workerPool.Start()
}

// Description:
// Stop a worker pool from the GoworkScheduler
func (p *GoworkScheduler) Stop() {
	log.Println("Stopping Go Work scheduler")
	p.workerPool.Stop()
}
