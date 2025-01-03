package service

import (
    "context"
    "log"
    "strconv"
    "time"

    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
)

type GameService struct {
    sqlite   *gorm.DB
    redis    *rueidis.Client
    GameRepo repository.GameRepositoryContract
}

func NewGameService(sqliteClient *gorm.DB, redisClient *rueidis.Client, repo repository.GameRepositoryContract) *GameService {
    return &GameService{
        sqlite:   sqliteClient,
        redis:    redisClient,
        GameRepo: repo,
    }
}

func (s *GameService) Matchmaking(ctx context.Context, userId, gameID string, distributedCards map[int][]string, kingCards, king string) {
    time.Sleep(1 * time.Second)
    lastMoveTimestamp := strconv.FormatInt(time.Now().Unix(), 10)
    s.GameRepo.Matchmaking(ctx, distributedCards, userId, gameID, lastMoveTimestamp, king, kingCards)
}

func (s *GameService) RemovePlayerFromWaitingList(username string) {
    s.GameRepo.RemovePlayerFromWaitingList(context.Background(), "matchmaking", username)
}

func (s *GameService) Subscribe(ctx context.Context, channel string, message func(rueidis.PubSubMessage)) error {
    err := (*s.redis).Receive(ctx, (*s.redis).B().Subscribe().Channel("game_creation").Build(), message)

    if err != nil {
        log.Printf("Error in subscribing to %v channel: %v", channel, err)
        return err
    }

    return nil
}
