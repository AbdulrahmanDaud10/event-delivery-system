package main

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

func ScheduleRetry(event FailedEvent) {
	// Increase the retry count
	event.RetryCount++

	// if retries are exhausted, log and potentially alert
	if event.RetryCount > MaxRetries {
		log.Printf("Failed to deliver event after %d attempts: %v", MaxRetries, event.Event)
		// NotifyAdmin("Event Delivery Failed", fmt.Printf("Failed to deliver event after %d attempts: %v", MaxRetries, event.Event))
		return
	}

	// Calculate next retry time with exponential backoff
	backOffDuration := time.Duration(math.Pow(2, float64(event.RetryCount))) * time.Second
	retryTimestamp := time.Now().Add(backOffDuration).Unix()

	eventJSON, _ := json.Marshal(event)
	rdb.ZAdd(ctx, "retry_events", &redis.Z{
		Score:  float64(retryTimestamp),
		Member: eventJSON,
	})
}
