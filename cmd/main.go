package main

import (
	"context"
	"log"
	"tender-backend/internal/pkg/config"
	"tender-backend/internal/http"
	db "tender-backend/internal/usecase/postgres"
	"tender-backend/internal/http/handlers"

	"github.com/redis/go-redis/v9" // Correct Redis import for v9
)

var redisClient *redis.Client

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	db.ConnectDB()
	defer db.CloseDB()

	// Initialize Redis
	InitRedis()
	defer redisClient.Close()

	// Initialize HTTP handlers
	h := handlers.NewHttpHandler(db.DB, redisClient)

	// Create and run the router
	r := http.NewGinRouter(h)
	err := r.Run(config.GlobalConfig.AppPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// InitRedis initializes the Redis client connection.
func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.GlobalConfig.Redis.RedisAddr, // Redis address from config
		Password: config.GlobalConfig.Redis.RedisPass, // Redis password from config
		DB:       0,                                   // Default DB
	})

	// Test the connection
	_, err := redisClient.Ping(context.Background()).Result() // Added context argument
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}
