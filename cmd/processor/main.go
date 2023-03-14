package main

import (
	"time"
	"twitch_chat_analysis/pkg/processing"
	"twitch_chat_analysis/pkg/queue"

	"github.com/Wr4thon/gorabbitmq"
	"github.com/go-redis/redis"
)

func main() {
	// setup redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379", //your redis address with the port. Default is localhost:6379
		Password:    "",
		DB:          0,
		ReadTimeout: 5 * time.Second,
	})

	processingService := processing.New(redisClient)

	// setup queue for consuming
	queue, err := queue.InitializeQueue()
	if err != nil {
		panic(err)
	}

	consumerSettings := gorabbitmq.ConsumerSettings{AutoAck: false, Exclusive: false, NoLocal: false, NoWait: false}

	deliveryConsumer := gorabbitmq.DeliveryConsumer(processingService.Consume)
	err = queue.RegisterConsumer(consumerSettings, deliveryConsumer)
	if err != nil {
		panic(err)
	}

	queue.Consume(consumerSettings)

}
