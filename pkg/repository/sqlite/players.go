package sqliteRepo

import (
	"github.com/alirezadp10/hokm/pkg/api/request"
	"github.com/alirezadp10/hokm/pkg/model"
	"github.com/alirezadp10/hokm/pkg/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

var _ repository.PlayersRepositoryContract = &PlayersRepository{}

type PlayersRepository struct {
	sqlite gorm.DB
}

func NewPlayersRepository(sqliteClient *gorm.DB) *PlayersRepository {
	return &PlayersRepository{
		sqlite: *sqliteClient,
	}
}

func (r *PlayersRepository) CheckPlayerExistence(username string) bool {
	var count int64

	err := r.sqlite.Table("players").Where("username = ?", username).Count(&count).Error

	if err != nil {
		log.Fatal(err)
		return false
	}

	return count > 0
}

func (r *PlayersRepository) SavePlayer(user request.User, chatId int64) (*model.Player, error) {
	newPlayer := model.Player{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		ChatId:    chatId,
		UpdatedAt: time.Now(),
		JoinedAt:  time.Now(),
	}

	err := r.sqlite.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&newPlayer).Error

	return &newPlayer, err
}

func (r *PlayersRepository) AddPlayerToGame(username, gameID string) (*model.Game, error) {
	var player model.Player
	r.sqlite.First(&player, "username = ?", username)

	newGame := model.Game{GameId: gameID, PlayerId: player.Id, CreatedAt: time.Now()}

	err := r.sqlite.Create(&newGame).Error

	return &newGame, err
}

func (r *PlayersRepository) DoesPlayerBelongToGame(username, gameID string) (bool, error) {
	var count int64

	err := r.sqlite.Table("players").
		Joins("inner join games on games.player_id = players.id").
		Where("players.username = ?", username).
		Where("games.game_id = ?", gameID).
		Count(&count).Error

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return count > 0, nil
}

func (r *PlayersRepository) DoesPlayerHaveAnyActiveGame(username string) (*string, bool) {
	var result struct{ GameId string }

	r.sqlite.Table("players").
		Select("games.game_id").
		Joins("inner join games on games.player_id = players.id").
		Where("players.username = ?", username).
		Where("games.finished_at is null").
		Scan(&result)

	if result.GameId != "" {
		return &result.GameId, true
	}

	return nil, false
}

func (r *PlayersRepository) HasGameFinished(gameID string) (bool, error) {
	var game model.Game

	r.sqlite.Table("games").Where("game_id = ?", gameID).First(&game)

	return game.FinishedAt != nil, nil
}
