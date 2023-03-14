package main

import (
	"net/http"
	"twitch_chat_analysis/pkg/queue"

	"github.com/gin-gonic/gin"
)

// Request is used to publish a message on the rabbitMQ.
type Request struct {
	Sender  string `json:"sender" binding:"required"`
	Message string `json:"message" binding:"required"`
}

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
		var request Request

		err := c.Bind(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = queue.Send(c, request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "enqueued")
	})

	println(r.Run().Error())
}
