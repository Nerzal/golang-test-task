package processing

import (
	"encoding/json"
	"twitch_chat_analysis/pkg/domain"

	"github.com/Noobygames/amqp"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

const redisKey = "twitch_chat_analysis"

// Service provides message processing capabilities.
type Service struct {
	redisClient *redis.Client
}

// New creates a new instance of processing service.
func New(redisClient *redis.Client) *Service {
	return &Service{
		redisClient: redisClient,
	}
}

func (s *Service) Consume(delivery amqp.Delivery) error {
	const errMessage = "could not consume message"

	var message domain.Message
	err := json.Unmarshal(delivery.Body, &message)
	if err != nil {
		nackErr := delivery.Nack(false, false)
		if nackErr != nil {
			return errors.Wrap(err, errMessage)
		}

		return errors.Wrap(err, errMessage)
	}

	err = s.handleMessage(message)
	if err != nil {
		nackErr := delivery.Nack(false, true) // failed to handle message, requeue
		if err != nil {
			return errors.Wrap(nackErr, errMessage)
		}

		return errors.Wrap(err, errMessage)
	}

	err = delivery.Ack(false)
	if err != nil {
		return errors.Wrap(err, errMessage)
	}

	return nil
}

func (s *Service) handleMessage(message domain.Message) error {
	const errMessage = "could not push to redis"

	cmd := s.redisClient.LPush(redisKey, message)
	_, err := cmd.Result()
	if err != nil {
		return errors.Wrap(err, errMessage)
	}

	return nil
}
