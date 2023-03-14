package main

import (
	"encoding/json"
	"net/http"
	"time"
	"twitch_chat_analysis/pkg/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const redisKey = "twitch_chat_analysis"

func main() {
	// setup redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379", //your redis address with the port. Default is localhost:6379
		Password:    "",
		DB:          0,
		ReadTimeout: 5 * time.Second,
	})

	r := gin.Default()

	r.GET("/message/list", func(c *gin.Context) {
		cmd := redisClient.LRange(redisKey, 0, -1)
		listEntries, err := cmd.Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []domain.Message

		for i := range listEntries {
			entry := listEntries[i]

			var message domain.Message
			err := json.Unmarshal([]byte(entry), &message)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			response = append(response, message)
		}

		c.JSON(http.StatusOK, response)
	})

	println(r.Run().Error())
}
