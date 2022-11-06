package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/raaaaaaaay86/clustered-cron-job/cronjob"
	"github.com/raaaaaaaay86/clustered-cron-job/cronjob/job"
	"github.com/raaaaaaaay86/clustered-cron-job/redis"
	"github.com/raaaaaaaay86/clustered-cron-job/redislock"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	redisManager := redis.NewRedisManager()
	redisLockManager := redislock.NewRedisLockManager(redisManager)
	greetingJob := job.NewGreetingJob()
	cronManager := cronjob.NewCronManager(redisLockManager, greetingJob)

	go cronManager.Start()

	r := gin.Default()

	r.GET("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	r.Run(addr)
}
