package background_jobs

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type GoworkScheduler struct {
	enqueuer   *work.Enqueuer
	workerPool *work.WorkerPool
}

type JobContext struct{}

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

func (p *GoworkScheduler) EnqueueJob(jobName string, jobArguments map[string]interface{}) error {
	_, err := p.enqueuer.Enqueue(jobName, jobArguments)
	return err
}

func (p *GoworkScheduler) RegisterCronJob(name, cronSpec string) {
	p.workerPool.PeriodicallyEnqueue(cronSpec, name)
}

func (p *GoworkScheduler) RegisterJob(jobName string, handler func(args map[string]interface{}) error) {
	p.workerPool.Job(jobName, func(job *work.Job) error {
		args := make(map[string]interface{})
		for key, val := range job.Args {
			args[key] = val
		}
		return handler(args)
	})

}

func (p *GoworkScheduler) Start() {
	log.Println("Starting Go Work scheduler")
	p.workerPool.Start()
}

func (p *GoworkScheduler) Stop() {
	log.Println("Stopping Go Work scheduler")
	p.workerPool.Stop()
}
