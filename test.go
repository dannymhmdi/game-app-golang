package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // No password set
		DB:       0,  // Use default DB
	})
	ctx := context.Background()
	err := rdb.Set(ctx, "bob", 26, 1*time.Second).Err()
	if err != nil {
		fmt.Println("err", err)

		return
	}
	value, gErr := rdb.Get(ctx, "bob").Result()
	if gErr != nil {
		fmt.Println("gErr", gErr)

		return
	}
	fmt.Println("value", value)
}
