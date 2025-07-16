package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"notificationservice/api-server/internal"
)

func main() {
	// Get configuration from environment variables
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	port := getEnv("PORT", "8082")

	// Setup Redis options
	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
		DB:   0,
	}

	// Setup Asynq client for enqueueing tasks
	client := asynq.NewClient(redisOpt)
	defer client.Close()

	// Setup Asynq inspector for task monitoring
	inspector := asynq.NewInspector(redisOpt)

	// Setup router
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		// Setup service in the context
		service := internal.NewService(client, inspector)
		c.Set(internal.ServiceContextKey, service)
	})

	// Register routes
	router.POST("/tasks", internal.SendNotification)
	router.GET("/tasks/pending", internal.GetPendingTasks)
	router.GET("/tasks/completed", internal.GetCompletedTasks)
	router.GET("/tasks/failed", internal.GetFailedTasks)
	router.PUT("/tasks/:id/send-now", internal.SendNow)
	router.PUT("/tasks/:id/cancel", internal.CancelTask)

	// Start server
	log.Printf("API server starting on port %s...\n", port)
	log.Printf("Redis connected at %s\n", redisAddr)
	log.Printf("Endpoints available: POST /tasks, GET /tasks/pending, GET /tasks/completed, GET /tasks/failed, GET /metrics")

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
