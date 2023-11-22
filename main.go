package main

import "context"

var (
	rdb RedisClientInterface   // Redis Client
	ctx = context.Background() // Global context Redis operation
)

func main() {
	InitializingRedis() // Initializing the redis client
	defer rdb.Close()   // Ensuring the redis client is closed when the program is exited
}
