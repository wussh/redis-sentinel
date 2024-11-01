package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "host.docker.internal:6379", // Address of the Redis server
	})

	// Set a key
	err := rdb.Set(ctx, "mykey", "hello, redis!", 0).Err()
	if err != nil {
		log.Fatalf("Could not set key: %v", err)
	}
	fmt.Println("Key set: mykey = hello, redis!")

	// Get the key
	val, err := rdb.Get(ctx, "mykey").Result()
	if err != nil {
		log.Fatalf("Could not get key: %v", err)
	}
	fmt.Printf("Key retrieved: mykey = %s\n", val)

	// Clean up and close the connection
	err = rdb.Close()
	if err != nil {
		log.Fatalf("Could not close the connection: %v", err)
	}
}
