package main

import (
	"encoding/json"
	"log"
	"time"
)

// Event represent the structure for the coming event data
type Event struct {
	UserID  string `json:"user_id"`
	PayLoad string `json:"payload"`
}

// FallEvent represents an event that has failed delivery, along with the count of retry attempts
type FallEvent struct {
	Event      Event
	RetryCount int
}

// ProcessEvent continously tries to fetch events from Redis and Sends them to their destinations
func ProcessEvent() {
	for {
		ProcessEvent()
	}
}

// ProcessSingleEvent tries to fetch a single event from Redis and sends it to its destination
func ProcessSingleEvent() {
	// pop an event from the front of Redis list (blocking until one is available)
	eventJSON, err := rdb.BLPop(ctx, 0*time.Second, "events").Result()
	if err != nil {
		log.Printf("Error fetching event from Redis: %v", err)
		return
	}

	if len(eventJSON) < 2 {
		log.Printf("Error: Unexpected BLPop result format")
		return
	}

	var event Event
	err = json.Unmarshal([]byte(eventJSON[1]), &event)
	if err != nil {
		log.Printf("Error Unmarshling event: %v", err)
		return
	}

	success := SendToDestination(event)
	if !success {
		log.Printf("Failed to deliver event: %v", event)
		// scheduleRetry(failedEvent{
		// 	Event:      event,
		// 	RetryCount: 1, // Intial retry attempt
		// })
	}
}
