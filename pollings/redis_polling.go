package pollings

import (
	"fmt"
	"time"

	"github.com/alfanzain/custom-agent-allocation/config"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisPolling struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisPolling(client *redis.Client) *RedisPolling {
	return &RedisPolling{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisPolling) StartRedisPolling() {
	key := config.REDIS_QUEUE_CUSTOMERS_KEY
	for {
		length, err := r.client.LLen(r.ctx, key).Result()
		if err != nil {
			fmt.Println("[Redis Polling] Error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if length == 0 {
			fmt.Println("[Redis Polling] List", key, "is empty. Skipping...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Println("[Redis Polling] List", key, "has", length, "element(s). Fetching all elements...")

		elements, err := r.client.LRange(r.ctx, key, 0, -1).Result()
		if err != nil {
			fmt.Println("[Redis Polling] Error fetching elements:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for i, element := range elements {
			fmt.Printf("[Redis Polling] Element %d: %s\n", i+1, element)
		}

		time.Sleep(2 * time.Second)
	}
}
