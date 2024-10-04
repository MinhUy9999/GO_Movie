package config

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
)

// ... existing code ...

func ConnectRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Địa chỉ Redis server
		Password: "",               // Mật khẩu nếu có
		DB:       0,                // Số database mặc định
	})

	// Kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Không thể kết nối đến Redis: %v", err)
		return err
	}

	log.Println("Kết nối Redis thành công")
	return nil
}
