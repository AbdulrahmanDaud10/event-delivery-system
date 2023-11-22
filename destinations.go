package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Destination interface represent any target to which we want to send our event
type Destination interface {
	Send(event Event) bool
}

// MockDestination1: A destination that successeds 80% of the time and fails 20% of the time
type MockDestination1 struct {
}

func (md *MockDestination1) Send(event Event) bool {
	randNumber := rand.Intn(100)
	if randNumber < 80 {
		fmt.Printf("MockDestination1 Successfully received event: %v\n", event)
		return true
	}

	fmt.Printf("Mockdestination1 failed to successfully receive event: %v\n", event)
	return false
}

// MockDestination2:  A destination that introduces random delays (up to 2 seconds)
type MockDestination2 struct {
}

func (md *MockDestination2) Send(event Event) bool {
	randDuration := time.Duration(rand.Intn(2000) * int(time.Millisecond))
	time.Sleep(randDuration)
	fmt.Printf("MockDestination2 has successfully received event: %v\n", event)

	return true
}

// MockDestination3: A destination that always succeeds and also logs the received event
type MockDestination3 struct {
}

func (md *MockDestination3) Send(event Event) bool {
	fmt.Printf("MockDestination3 has successfully received event: %v\n", event)
	return true
}

// Mock function to simulate sending the event to a destination
// Randomly returns success or failure
func SendToDestination(event Event) bool {
	destinations := []Destination{
		&MockDestination1{},
		&MockDestination2{},
		&MockDestination3{},
	}

	for _, destination := range destinations {
		success := destination.Send(event)
		if !success {
			return false // Event delivery failed for one of the destinations
		}
	}

	return true
}
