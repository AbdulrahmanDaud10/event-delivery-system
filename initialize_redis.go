package main

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// InitializingRedis Sets up a connection to the Redis Server
func InitializingRedis() {
	myRedis := "myRedis:9999" // Redis Server Address
	client := redis.NewClient(&redis.Options{
		Addr:     myRedis,
		Password: "",
		DB:       0,
	})

	fmt.Println("Initialize Redis", myRedis)
	rdb := &RedisClientWrapper{Client: client}

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
	}
}
