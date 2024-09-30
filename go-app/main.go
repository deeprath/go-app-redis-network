package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Define a struct for the table
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

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

	// Simulate the table data with 5 rows
	users := []User{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 1, Name: "Bob", Age: 30},
		{ID: 1, Name: "Charlie", Age: 22},
		{ID: 1, Name: "Diana", Age: 35},
		{ID: 1, Name: "Eve", Age: 28},
	}

	// Example function to store data in Redis
	err = cacheUsers(users, rdb)
	if err != nil {
		log.Fatalf("error caching users: %v", err)
	}

	// Retrieve and print the cached data from Redis
	cachedUsers, err := getCachedUsersWithPagination(rdb, 4, 1)
	if err != nil {
		log.Fatalf("error getting cached users: %v", err)
	}

	start := time.Now()
	// Print the retrieved data
	for _, user := range cachedUsers {
		fmt.Printf("User: %v, Age: %d\n", user.Name, user.Age)
	}
	timeElapsed := time.Since(start)
	fmt.Println("timeElasped : ", timeElapsed)

}

func cacheUsers(users []User, rdb *redis.Client) error {
	cacheKey := "users_table"

	// Convert users to JSON for storing in Redis
	usersJSON, err := json.Marshal(users)
	if err != nil {
		return err
	}

	// Store the users in Redis with TTL (e.g., 1 hour)
	err = rdb.Set(ctx, cacheKey, usersJSON, time.Hour).Err()
	if err != nil {
		return err
	}

	fmt.Println("Users cached successfully.")
	return nil
}

func getCachedUsersWithPagination(rdb *redis.Client, limit, offset int) ([]User, error) {
	cacheKey := "users_table"

	// Try fetching from Redis
	cachedUsers, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		fmt.Println("No users found in cache.")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Parse the cached data
	var users []User
	err = json.Unmarshal([]byte(cachedUsers), &users)
	if err != nil {
		return nil, err
	}

	// Handle pagination: limit and offset
	totalUsers := len(users)
	if offset > totalUsers {
		return nil, fmt.Errorf("offset out of range")
	}

	end := offset + limit
	if end > totalUsers {
		end = totalUsers
	}

	// Return the paginated slice of users
	return users[offset:end], nil
}
