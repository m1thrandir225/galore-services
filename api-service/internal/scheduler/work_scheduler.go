package scheduler

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// WorkScheduler implements the Service via the GoWork package
type WorkScheduler struct {
	enqueuer   *work.Enqueuer
	workerPool *work.WorkerPool
}

type JobContext struct{}

// NewWorkScheduler returns a GoWork Scheduler implementation.
func NewWorkScheduler(namespace string, redisUrl string) *WorkScheduler {
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
	return &WorkScheduler{
		enqueuer:   enqueuer,
		workerPool: workerPool,
	}
}

func (p *WorkScheduler) EnqueueJob(jobName string, jobArguments map[string]interface{}) error {
	_, err := p.enqueuer.Enqueue(jobName, jobArguments)
	return err
}

func (p *WorkScheduler) RegisterCronJob(name, cronSpec string) {
	p.workerPool.PeriodicallyEnqueue(cronSpec, name)
}

func (p *WorkScheduler) RegisterJob(jobName string, noRetry bool, handler func(args map[string]interface{}) error) {
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

func (p *WorkScheduler) Start() {
	log.Println("Starting Go Work scheduler")
	p.workerPool.Start()
}

func (p *WorkScheduler) Stop() {
	log.Println("Stopping Go Work scheduler")
	p.workerPool.Stop()
}
