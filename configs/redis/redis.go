package redis

import (
	"context"
	"errors"
	"os"
	"time"
	"github.com/go-redis/redis/v8"
)

// ConnectRedis establishes a connection to the Redis server
func ConnectRedis() (*redis.Client, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, errors.New("REDIS_URL environment variable is not set")
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
