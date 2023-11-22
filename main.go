package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	rdb    RedisClientInterface   // Redis Client
	ctx    = context.Background() // Global context Redis operation
	logger = logrus.New()         // logrus logger instance
)

const (
	MaxRetries = 10 // Max attempts to retry sending an event
)

func main() {
	InitializingRedis() // Initializing the redis client
	defer rdb.Close()   // Ensuring the redis client is closed when the program is exited

	// Start a Go routine for processing events.
	go ProcessEvent()

	// Start a separate routine to process failed events.
	for i := 0; i < 10; i++ {
		go ProcessFailedEvents(i)
	}

	// Set up our HTTP server.
	http.HandleFunc("/ingest", IngestEventHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
