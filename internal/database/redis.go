package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// This is the global variable our signup.go file is trying to use!
var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Points to the Docker container named "redis"
		Password: "",           // No password by default
		DB:       0,            // Default DB
	})

	// Send a quick "Ping" to see if Redis is actually awake
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("❌ REDIS ERROR: Failed to connect to Redis container!", err)
	} else {
		fmt.Println("✅ REDIS CONNECTED: Ready to store OTPs!")
	}
}