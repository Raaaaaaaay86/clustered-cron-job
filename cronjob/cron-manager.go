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
	lock := c.obtainLock("your_app_name") // Thr process will be blocked at this line until obtain the lock.

	go c.extendLock(lock) // Open new goroutine to extend lock ownership repeatedly.

	c.Cron.AddJob("@every 1s", c.GreetingJob)
	c.Cron.Start()
}

func (c CronManager) obtainLock(name string) *redsync.Mutex {
	rs := c.RedisLockManager.Client()

	for {
		// Retry every N seconds until obtain the lock.
		time.Sleep(10 * time.Second)

		lock := rs.NewMutex(name, redsync.WithExpiry(30*time.Second))

		if err := lock.Lock(); err != nil {
			c.logger.Printf("[CronManager] Obtained the lock failed: %s", err)

		} else {
			c.logger.Printf("[CronManager] Obtained the lock successfully.")

			return lock
		}
	}
}

func (c CronManager) extendLock(lock *redsync.Mutex) {
	for {
		// Refresh lock TTL every N seconds.
		time.Sleep(10 * time.Second)

		ok, err := lock.Extend()
		if err != nil {
			c.logger.Printf("[CronManager] Got error when extending lock (%s): %s", lock.Name(), err.Error())

			break
		}

		if !ok {
			c.logger.Printf("[CronManager] Extend lock failed (%s)", lock.Name())

			break
		}

		c.logger.Printf("[CronManager] Extend lock succeed (%s)", lock.Name())
	}
}
