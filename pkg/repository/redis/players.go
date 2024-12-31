package redisRepo

import (
    "github.com/alirezadp10/hokm/pkg/api/request"
    "github.com/alirezadp10/hokm/pkg/model"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/redis/rueidis"
)

var _ repository.PlayersRepositoryContract = &PlayersRepository{}

type PlayersRepository struct {
    redis rueidis.Client
}

func (p PlayersRepository) CheckPlayerExistence(username string) bool {
    //TODO implement me
    panic("implement me")
}

func (p PlayersRepository) SavePlayer(user request.User, chatId int64) (*model.Player, error) {
    //TODO implement me
    panic("implement me")
}

func (p PlayersRepository) AddPlayerToGame(username, gameID string) (*model.Game, error) {
    //TODO implement me
    panic("implement me")
}

func NewPlayersRepository(redisClient *rueidis.Client) *PlayersRepository {
    return &PlayersRepository{
        redis: *redisClient,
    }
}
