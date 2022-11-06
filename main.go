package main

import "github.com/gin-gonic/gin"

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	redisManager := redis.NewRedisManager()
	redisLockManager := redislock.NewRedisLockManager(redisManager)
	greetingJob := job.NewGreetingJob()
	cronManager := cronjob.NewCronManager(redisLockManager, greetingJob)

	r := gin.Default()

	r.GET("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	go cronManager.Start()

}
