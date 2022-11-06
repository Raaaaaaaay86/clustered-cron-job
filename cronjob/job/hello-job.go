package job

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type IGreetingJob interface {
	cron.Job
}

type GreetingJob struct {
	logger *log.Logger
}

func NewGreetingJob() IGreetingJob {
	return GreetingJob{
		logger: log.Default(),
	}
}

func (g GreetingJob) Run() {
	g.logger.Printf("[GreetingJob] Greeting. Current time is %s", time.Now())
}
