package queue

import (
	"github.com/Wr4thon/gorabbitmq"
	"github.com/pkg/errors"
)

const (
	queueName        = "message-queue" // name for the queue
	rabbitMQUser     = "user"          // name from docker-compose.yml
	rabbitMQPassword = "password"      // pw from docker-compose.yml
	rabbitMQHost     = "localhost"     // change when not running on localhost
	rabbitMQPort     = 7001            // port from docker-compose.yml
)

// InitializeQueue initializes a new queue connection.
// After calling this function, the queue will be ready to publish new messages or subscribe to a queue.
func InitializeQueue() (gorabbitmq.Queue, error) {
	connectionSettings := gorabbitmq.ConnectionSettings{
		UserName: rabbitMQUser,
		Password: rabbitMQPassword,
		Host:     rabbitMQHost,
		Port:     rabbitMQPort,
	}

	// use default values
	channelSettings := gorabbitmq.ChannelSettings{}

	// create a new connection
	qConnector, err := gorabbitmq.NewConnection(connectionSettings, channelSettings)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to new queue")
	}

	queueSettings := gorabbitmq.QueueSettings{
		QueueName:        queueName,
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
