package redis

import (
	"log"
	"os"

	"github.com/go-redis/redis/v9"
)

type IRedisManager interface {
	Client() *redis.Client
}

type RedisManager struct {
	logger *log.Logger
}

func NewRedisManager() IRedisManager {
	return RedisManager{
		logger: log.Default(),
	}
}

var redisClientInstance *redis.Client

func (r RedisManager) Client() *redis.Client {
	if redisClientInstance == nil {
		return r.createConnection()
	}

	return redisClientInstance
}

func (r RedisManager) createConnection() *redis.Client {
	address := os.Getenv("REDIS_ADDRESS")
	redisClientInstance = redis.NewClient(&redis.Options{
		Addr: address,
	})

	r.logger.Printf("[RedisManager] Redis connection established. (%s)", address)

	return redisClientInstance
}
