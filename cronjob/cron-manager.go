package cronjob

import (
	"log"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/raaaaaaaay86/clustered-cron-job/cronjob/job"
	"github.com/raaaaaaaay86/clustered-cron-job/redislock"
	"github.com/robfig/cron/v3"
)

type ICronManager interface {
	Start()
}

type CronManager struct {
	RedisLockManager redislock.IRedisLockManager
	GreetingJob      job.IGreetingJob
	Cron             *cron.Cron
	logger           log.Logger
}

type JobSpec struct {
	TimeSpec string
	Job      cron.Job
}

func NewCronManager(
	redisLockManager redislock.IRedisLockManager,
	greetingJob job.IGreetingJob,
) ICronManager {
	return CronManager{
		Cron:             cron.New(),
		RedisLockManager: redisLockManager,
		GreetingJob:      greetingJob,
		logger:           *log.Default(),
	}
}

func (c CronManager) Start() {
	for {
		lock := c.obtainLock("your_app_name") // Thr process will be blocked at this line until obtain the lock.

		stopChan := make(chan int)

		go c.extendLock(lock, stopChan) // Open new goroutine to extend lock ownership repeatedly.

		job := []JobSpec{
			{
				TimeSpec: "@every 1s",
				Job:      c.GreetingJob,
			},
		}

		entryIds := c.addJobs(job)
		c.Cron.Start()

		<-stopChan // Block here to wait stop signal

		c.removeJobs(entryIds)
		c.Cron.Stop()
	}
}

func (c CronManager) obtainLock(name string) *redsync.Mutex {
	// rs := c.RedisLockManager.Client()

	for {
		// Retry every N seconds until obtain the lock.
		time.Sleep(10 * time.Second)

		lock := c.RedisLockManager.NewLock(name, 30*time.Second)

		if err := lock.Lock(); err != nil {
			c.logger.Printf("[CronManager] Obtained the lock failed: %s", err)

		} else {
			c.logger.Printf("[CronManager] Obtained the lock successfully.")

			return lock
		}
	}
}

func (c CronManager) extendLock(lock *redsync.Mutex, stopChan chan<- int) {
	for {
		// Refresh lock TTL every N seconds.
		time.Sleep(10 * time.Second)

		ok, err := lock.Extend()
		if err != nil {
			c.logger.Printf("[CronManager] Got error when extending lock (%s): %s", lock.Name(), err.Error())

			stopChan <- 1
			break
		}

		if !ok {
			c.logger.Printf("[CronManager] Extend lock failed (%s)", lock.Name())

			stopChan <- 1
			break
		}

		c.logger.Printf("[CronManager] Extend lock succeed (%s)", lock.Name())
	}
}

func (c CronManager) addJobs(jobSpecs []JobSpec) []cron.EntryID {
	entryIds := make([]cron.EntryID, 0)

	for _, jobSpec := range jobSpecs {
		entryId, err := c.Cron.AddJob(jobSpec.TimeSpec, jobSpec.Job)
		if err != nil {
			panic(err)
		}

		entryIds = append(entryIds, entryId)
	}

	return entryIds
}

func (c CronManager) removeJobs(entryIds []cron.EntryID) {
	for _, entryId := range entryIds {
		c.Cron.Remove(entryId)
	}
}
