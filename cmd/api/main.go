package main

import (
	"net/http"
	"twitch_chat_analysis/pkg/domain"
	"twitch_chat_analysis/pkg/queue"

	"github.com/gin-gonic/gin"
)

func main() {
	queue, err := queue.InitializeQueue()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "worked")
	})

	r.POST("/message", func(c *gin.Context) {
		var payload domain.Message

		err := c.Bind(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = queue.Send(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "enqueued")
	})

	println(r.Run().Error())
}
