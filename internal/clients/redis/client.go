package redis

import (
	"fmt"

	cfg "github.com/fiap-challenger-soat/hackthon-soat-process-worker/config"

	"github.com/go-redis/redis"
)

func NewRedisClient() (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     cfg.Vars.RedisAddress,
		// Password: cfg.Vars.RedisPassword,
		DB:       cfg.Vars.RedisDB,
	}

	rdb := redis.NewClient(opts)

	if err := rdb.Ping().Err(); err != nil {
		return nil, fmt.Errorf("failed to connect and ping Redis: %w", err)
	}

	return rdb, nil
}
