package services

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type QueueService struct {
	RedisClient *redis.Client
	Ctx         context.Context
}

func NewQueueService(redisClient *redis.Client, ctx context.Context) *QueueService {
	return &QueueService{
		RedisClient: redisClient,
		Ctx:         ctx,
	}
}

func (s *QueueService) EnqueueCustomer(queueName string, roomID string) error {
	err := s.RedisClient.RPush(s.Ctx, queueName, roomID).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue customer: %w", err)
	}
	return nil
}

func (s *QueueService) DequeueCustomer(queueName string) (string, error) {
	val, err := s.RedisClient.LPop(s.Ctx, queueName).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("queue is empty")
	}
	if err != nil {
		return "", fmt.Errorf("failed to dequeue customer: %w", err)
	}

	return val, nil
}

func (s *QueueService) DoesQueueExists(queueName string) (bool, error) {
	length, err := s.RedisClient.LLen(s.Ctx, queueName).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check queue length: %w", err)
	}

	return length > 0, nil
}
