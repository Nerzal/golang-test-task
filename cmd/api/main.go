package main

import (
	"net/http"

	"github.com/Clarilab/gorabbitmq/v4"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Request is used to publish a message on the rabbitMQ.
type Request struct {
	Sender  string `json:"sender" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func main() {
	queue, err := initializeQueue()
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

const (
	queueName        = "message-queue" // name for the queue
	rabbitMQUser     = "user"          // name from docker-compose.yml
	rabbitMQPassword = "password"      // pw from docker-compose.yml
	rabbitMQHost     = "localhost"     // change when not running on localhost
	rabbitMQPort     = 7001            // port from docker-compose.yml
)

func initializeQueue() (gorabbitmq.Queue, error) {
	connectionSettings := gorabbitmq.ConnectionSettings{
		UserName: rabbitMQUser,
		Password: rabbitMQPassword,
		Host:     rabbitMQHost,
		Port:     rabbitMQPort,
	}

	channelSettings := gorabbitmq.ChannelSettings{
		UsePrefetch:   false,
		PrefetchCount: 0,
	}

	qConnector, err := gorabbitmq.NewConnection(connectionSettings, channelSettings)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to new queue")
	}

	queueSettings := gorabbitmq.QueueSettings{
		QueueName:        queueName, // "ExampleQueueName"
		Durable:          true,
		DeleteWhenUnused: false,
		Exclusive:        false,
		NoWait:           false,
	}

	q, err := qConnector.ConnectToQueue(queueSettings)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to queue")
	}

	return q, nil
}
