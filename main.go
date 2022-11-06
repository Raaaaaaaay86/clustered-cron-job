package main

import "github.com/gin-gonic/gin"

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
