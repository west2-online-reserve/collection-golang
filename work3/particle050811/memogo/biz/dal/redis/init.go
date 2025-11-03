package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// RDB 全局 Redis 客户端实例
	RDB *redis.Client
)

// Init 初始化 Redis 连接
func Init() {
	// 从环境变量读取配置，提供默认值
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")
	db := getEnvAsInt("REDIS_DB", 0)

	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v (caching will be disabled)", err)
		RDB = nil
		return
	}

	log.Println("✅ Redis connected successfully")
}

// Close 关闭 Redis 连接
func Close() {
	if RDB != nil {
		if err := RDB.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取整数类型的环境变量
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
