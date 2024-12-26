package sqlite

import (
    "github.com/alirezadp10/hokm/internal/database"
    "github.com/alirezadp10/hokm/internal/request"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "log"
    "time"
)

func SavePlayer(db *gorm.DB, user request.User, chatId int64) (*database.Player, error) {
    newPlayer := database.Player{
        Id:        user.ID,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Username:  user.Username,
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

func AddPlayerToGame(db *gorm.DB, username, gameID string) (*database.Game, error) {
    var player database.Player
    db.First(&player, "username = ?", username)

    newGame := database.Game{GameId: gameID, PlayerId: player.Id, CreatedAt: time.Now()}

    err := db.Create(&newGame).Error

    return &newGame, err
}

func CheckPlayerExistence(db *gorm.DB, username string) bool {
    var count int64

    err := db.Table("players").Where("username = ?", username).Count(&count).Error

    if err != nil {
        log.Fatal(err)
        return false
    }

    return count > 0
}
