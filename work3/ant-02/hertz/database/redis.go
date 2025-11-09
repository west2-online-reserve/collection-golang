package database

import (
	"context"
	"fmt"
	"hertz/config"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
	ctx    context.Context
}

type RedisClient interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	RPush(key string, value interface{}) error
	LRange(key string, first int64, second int64) ([]string, error)
	Del(key []string) error
	Exists(key string) (bool, error)
	LLen(key string) (int64, error)
}

var (
	redisOnce sync.Once
	instance  *redisClient
)

func GetRedis() RedisClient {
	redisOnce.Do(func() {
		cfg := config.GetConfig()
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Db,
		})
		ctx := context.Background()
		_, err := client.Ping(ctx).Result()
		if err != nil {
			panic(fmt.Sprintf("无法连接到Redis: %v", err))
		}
		instance = &redisClient{
			client: client,
			ctx:    context.Background(),
		}
	})
	return instance
}

func (r *redisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *redisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *redisClient) RPush(key string, value interface{}) error {
	_, err := r.client.RPush(r.ctx, key, value).Result()
	return err
}

func (r *redisClient) LRange(key string, first int64, second int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, first, second).Result()
}

func (r *redisClient) Del(key []string) error {
	_, err := r.client.Del(r.ctx, key...).Result()
	return err
}

func (r *redisClient) Exists(key string) (bool, error) {
	cnt, err := r.client.Exists(r.ctx, key).Result()
	return cnt == 1, err
}

func (r *redisClient) LLen(key string) (int64, error) {
	return r.client.LLen(r.ctx, key).Result()
}
