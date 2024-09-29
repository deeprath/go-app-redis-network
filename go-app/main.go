package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Get Redis connection details from environment variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// redisAddr := "redis-service:6379" // redis-service is the name of the Redis container

	// Set up Redis client with username and password
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // Redis address (host:port)
		Username: redisUsername, // Redis username
		Password: redisPassword, // Redis password
		DB:       0,             // Default database
	})

	// Test connection with a Ping command
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		fmt.Println("Successfully connected to Redis!")
	}

	// Set a test key
	err = rdb.Set(ctx, "test_key", "Hello from Go with Auth!", 0).Err()
	if err != nil {
		log.Fatalf("Could not set key: %v", err)
	}

	// Get the test key
	val, err := rdb.Get(ctx, "test_key").Result()
	if err != nil {
		log.Fatalf("Could not get key: %v", err)
	}

	fmt.Printf("test_key: %s\n", val)
}
