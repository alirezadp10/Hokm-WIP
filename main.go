package main

import (
    "context"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/hokm"
    "github.com/google/uuid"
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
    gameId := uuid.New().String()
    distributedCards := hokm.DistributeCards()
    redis.Matchmaking(ctx, client, "1", gameId, distributedCards)
    redis.Matchmaking(ctx, client, "2", gameId, distributedCards)
    redis.Matchmaking(ctx, client, "3", gameId, distributedCards)
    redis.Matchmaking(ctx, client, "4", gameId, distributedCards)
}
