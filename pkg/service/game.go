package service

import (
	"context"
	"github.com/alirezadp10/hokm/pkg/repository"
	"github.com/redis/rueidis"
	"gorm.io/gorm"
	"log"
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

func (s *GameService) Matchmaking(ctx context.Context, userId string) {
	s.GameRepo.Matchmaking(ctx, userId)
}

func (s *GameService) CheckAnyExistingGameForPlayer(ctx context.Context, username string) (string, bool) {
	var gameID string
	var exists bool
	gameID, exists, err := s.GameRepo.CheckAnyExistingGameForPlayer(ctx, username)
	if err != nil {
		log.Printf("Error checking existing game for player %v: %v", username, err)
		return "", false
	}
	return gameID, exists
}

func (s *GameService) AddPlayerToWaitingList(ctx context.Context, username string) {
	s.GameRepo.AddPlayerToWaitingList(ctx, username)
}

func (s *GameService) RemovePlayerFromWaitingList(ctx context.Context, username string) {
	s.GameRepo.RemovePlayerFromWaitingList(ctx, username)
}

func (s *GameService) Subscribe(ctx context.Context, channel string, message func(rueidis.PubSubMessage)) error {
	err := (*s.redis).Receive(ctx, (*s.redis).B().Subscribe().Channel("game_creation").Build(), message)

	if err != nil {
		log.Printf("Error in subscribing to %v channel: %v", channel, err)
		return err
	}

	return nil
}
