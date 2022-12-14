package redislock

import (
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/raaaaaaaay86/clustered-cron-job/redis"
)

type IRedisLockManager interface {
	Client() *redsync.Redsync
	NewLock(name string, expiry time.Duration) *redsync.Mutex
}

type RedisLockManager struct {
	RedisManager redis.IRedisManager
}

var redsyncInstance *redsync.Redsync

func NewRedisLockManager(redisManager redis.IRedisManager) IRedisLockManager {
	return RedisLockManager{
		RedisManager: redisManager,
	}
}

func (r RedisLockManager) Client() *redsync.Redsync {
	if redsyncInstance == nil {
		pool := goredis.NewPool(r.RedisManager.Client())
		redsyncInstance = redsync.New(pool)
	}

	return redsyncInstance
}

func (r RedisLockManager) NewLock(name string, expiry time.Duration) *redsync.Mutex {
	return r.Client().NewMutex(name, redsync.WithExpiry(expiry))
}
