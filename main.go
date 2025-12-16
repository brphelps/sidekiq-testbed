package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jrallison/go-workers"
)

// SleepWorker is a worker that sleeps for 10 seconds
func SleepWorker(message *workers.Msg) {
	log.Printf("Starting job at %s", time.Now().Format(time.RFC3339))
	time.Sleep(10 * time.Second)
	log.Printf("Completed job at %s", time.Now().Format(time.RFC3339))
}

func scheduleJobs() {
	// Schedule job to run twice per minute (every 30 seconds)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Enqueue initial job immediately
	_, err := workers.Enqueue("sleep_queue", "SleepWorker", []interface{}{})
	if err != nil {
		log.Printf("Error enqueueing initial job: %v", err)
	} else {
		log.Println("Initial job enqueued")
	}

	for {
		select {
		case <-ticker.C:
			_, err := workers.Enqueue("sleep_queue", "SleepWorker", []interface{}{})
			if err != nil {
				log.Printf("Error enqueueing job: %v", err)
			} else {
				log.Println("Job enqueued")
			}
		}
	}
}

func main() {
	// Get Redis connection details from environment variables
	redisServer := os.Getenv("REDIS_SERVER")
	if redisServer == "" {
		redisServer = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		redisDB = "0"
	}

	// Configure workers
	workers.Configure(map[string]string{
		"server":   redisServer,
		"password": redisPassword,
		"database": redisDB,
		"pool":     "10",
		"process":  "1",
	})

	// Register the worker
	workers.Process("sleep_queue", SleepWorker, 1)

	log.Printf("Starting Sidekiq worker connected to Redis at %s", redisServer)

	// Start the scheduler in a separate goroutine
	go scheduleJobs()

	// Start processing jobs
	go workers.Run()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down gracefully...")
	workers.Quit()
}
