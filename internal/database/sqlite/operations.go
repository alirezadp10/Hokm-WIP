package sqlite

import (
    "github.com/alirezadp10/hokm/internal/database"
    "gopkg.in/telebot.v4"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "time"
)

func SavePlayer(db *gorm.DB, player *telebot.User, chatId int64) (*database.Player, error) {
    newPlayer := database.Player{
        Id:        player.ID,
        FirstName: player.FirstName,
        LastName:  player.LastName,
        Username:  player.Username,
        ChatId:    chatId,
        UpdatedAt: time.Now(),
        JoinedAt:  time.Now(),
    }

    err := db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        DoNothing: true,
    }).Create(&newPlayer).Error

    return &newPlayer, err
}

func AddPlayerToGame(db *gorm.DB, playerId int64, gameId string) (*database.Game, error) {
    newGame := database.Game{
        GameId:    gameId,
        PlayerId:  playerId,
        CreatedAt: time.Now(),
    }

    err := db.Create(&newGame).Error

    return &newGame, err
}

func DoesPlayerHaveAnActiveGame(db *gorm.DB, playerId int64) (*string, bool) {
    return nil, false
}
