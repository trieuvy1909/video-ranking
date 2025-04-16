package redis

import (
	"context"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"github.com/go-redis/redis/v8"
)

// ConnectRedis establishes a connection to the Redis server
func ConnectRedis() (*redis.Client, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Printf("REDIS_URL not set")
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		return nil, err
	}

	log.Println("Connected to Redis successfully")
	return client, nil
}
