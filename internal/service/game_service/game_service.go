package game_service

import (
    "context"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/service/cards_service"
    "github.com/redis/rueidis"
    "strconv"
    "time"
)

// Matchmaking assigns players to a game and initializes game data in Redis
func Matchmaking(ctx context.Context, client rueidis.Client, userId, gameId string) {
    time.Sleep(1 * time.Second)                                   // Simulate delay for matchmaking
    distributedCards := cards_service.DistributeCards()           // Distribute cards among players
    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10) // Record timestamp
    kingCards, king := cards_service.ChooseFirstKing()            // Determine king cards and player
    redis.Matchmaking(ctx, client, distributedCards, userId, gameId, lastMoveTimestamp, king, kingCards)
}
