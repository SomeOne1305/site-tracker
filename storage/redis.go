package storage

import (
	"context"
	"visit-tracker/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	options, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		panic("Failed to parse Redis URL: " + err.Error())
	}

	client := redis.NewClient(options)

	// Test the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return &RedisClient{Client: client}
}
