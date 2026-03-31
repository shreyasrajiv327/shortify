package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Increment(key string) (int64, error) {
	return r.client.Incr(context.Background(), key).Result()
}

func (r *RedisClient) Expire(key string, duration time.Duration) error {
	return r.client.Expire(context.Background(), key, duration).Err()
}

func (r *RedisClient) Set(key string, value string) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}