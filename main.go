package main

import (
	"context"

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
}
