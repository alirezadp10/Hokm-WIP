package main

import (
    "context"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/joho/godotenv"
    "github.com/redis/rueidis"
    "log"
)

func main() {
    _ = godotenv.Load()

    client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
        log.Fatalf("could not connect to Redis: %v", err)
    }
    defer client.Close()

    ctx := context.Background()

    err = redis.SetTrump(ctx, client, "4", "b3")
    if err != nil {
        return
    }
}
