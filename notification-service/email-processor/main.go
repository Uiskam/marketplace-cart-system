package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"notificationservice/email-processor/internal"
)

func main() {
	// Get configuration from environment variables
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")

	// Setup srv
	redisOpt := asynq.RedisClientOpt{
		Addr: redisAddr,
		DB:   0,
	}

	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 1,
			Queues: map[string]int{
				"email_high": 9,
				"email_low":  1,
			},
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Duration(n) * time.Second
			},
		},
	)

	// Register email processor

	mux := asynq.NewServeMux()
	mux.Handle(internal.TaskEmail, internal.NewProcessor())

	// Setup signal handling for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Start the srv
	log.Println("Email processor starting...")
	if err := srv.Start(mux); err != nil {
		log.Fatalf("Failed to start srv: %v", err)
	}

	// Wait for termination signal
	<-done
	log.Println("Email processor shutting down...")
	srv.Shutdown()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
