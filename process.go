package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Event represent the structure for the coming event data
type Event struct {
	UserID  string `json:"user_id"`
	PayLoad string `json:"payload"`
}

// FailedEvent represents an event that has failed delivery, along with the count of retry attempts
type FailedEvent struct {
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
		ScheduleRetry(FailedEvent{
			Event:      event,
			RetryCount: 1, // Intial retry attempt
		})
	}
}

func ProcessFailedEvents(workerID int) {
	for {
		now := time.Now().Unix()

		// fetch events scheduled for a retry up to the current timestamp
		events, err := rdb.ZRangeByScoreWithScores(ctx, "retry_events", &redis.ZRangeBy{
			Min:    "0",
			Max:    fmt.Sprintf("%d", now),
			Offset: 0,
			Count:  1,
		}).Result()

		if err != nil {
			logger.WithFields(logrus.Fields{
				"workerID": workerID,
				"error":    err,
			}).Error("Failed to fetch events for retry")
			time.Sleep(10 * time.Second)
			continue
		}

		if len(events) == 0 {
			time.Sleep(10 * time.Second) // No events ready for retry, sleep for a while
			continue
		}

		var failedEvent FailedEvent
		err = json.Unmarshal([]byte(events[0].Member.(string)), &failedEvent)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"workerID": workerID,
				"error":    failedEvent.Event,
			}).Error("Error unmarshalling event failed")
			continue
		}

		success := SendToDestination(failedEvent.Event)
		if !success {
			logger.WithFields(logrus.Fields{
				"workerID": workerID,
				"error":    failedEvent.Event,
			}).Info("Event delivered successfuly on retry")

			rdb.ZRem(ctx, "retry_events", events[0].Member)
		} else {
			if failedEvent.RetryCount > MaxRetries {
				logger.WithFields(logrus.Fields{
					"workerID": workerID,
					"error":    failedEvent.Event,
				}).Warn("Maximum  retries exhausted  for event")
			} else {
				logger.WithFields(logrus.Fields{
					"workerID": workerID,
					"error":    failedEvent.Event,
				}).Info("Retry failed for event, rescheduling")
			}

			rdb.ZRem(ctx, "retry_events", events[0].Member)
			ScheduleRetry(failedEvent)
		}
	}
}
