package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the global variable we will use in our controllers
var RedisClient *redis.Client

func ConnectRedis() error {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" || port == "" {
		return fmt.Errorf("CRITICAL: Missing Redis environment variables in .env")
	}

	address := fmt.Sprintf("%s:%s", host, port)

	// Initialize the Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, 
		DB:       0,        // Use default DB
	})

	// Ping the server to verify the connection is actually alive
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	// Assign the active connection to our global variable
	RedisClient = client
	log.Println("✅ REDIS CONNECTED: Ready to store OTPs!")
	
	return nil
}