package processing

import (
	"context"
	"encoding/json"
	"time"
	"twitch_chat_analysis/pkg/domain"

	"github.com/Noobygames/amqp"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const redisKey = "twitch_chat_analysis_sorted_set"

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

	ctx := context.TODO()

	err = s.handleMessage(ctx, message)
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

func (s *Service) handleMessage(ctx context.Context, message domain.Message) error {
	const errMessage = "could not push to redis"

	cmd := s.redisClient.ZAdd(
		ctx, redisKey,
		redis.Z{Score: float64(time.Now().Unix()), Member: message},
	)

	_, err := cmd.Result()
	if err != nil {
		return errors.Wrap(err, errMessage)
	}

	return nil
}
