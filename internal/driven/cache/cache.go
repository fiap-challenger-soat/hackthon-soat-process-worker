package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const expiration = 24 * time.Hour // Set expiration time for job status cache

type JobCache interface {
	SetJobStatus(jobID, status string) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) SetJobStatus(jobID, status string) error {
	key := fmt.Sprintf("job_status:%s", jobID)
	if err := c.client.Set(key, status, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set job status for %s in Redis: %w", jobID, err)
	}
	return nil
}
