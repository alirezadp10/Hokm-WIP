package service

import (
    "context"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/repository"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
    "strconv"
    "time"
)

type GameService struct {
    sqlite   *gorm.DB
    redis    rueidis.Client
    GameRepo repository.GameRepository
}

func NewGameService(repo repository.GameRepository, sqlite *gorm.DB, redis rueidis.Client) *GameService {
    return &GameService{
        GameRepo: repo,
        sqlite:   sqlite,
        redis:    redis,
    }
}

func (s *GameService) Matchmaking(ctx context.Context, client rueidis.Client, userId, gameID string, distributedCards []string, king, kingCards string) {
    time.Sleep(1 * time.Second)
    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
    redis.Matchmaking(ctx, client, distributedCards, userId, gameID, lastMoveTimestamp, king, kingCards)
}
