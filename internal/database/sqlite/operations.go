package sqlite

import (
    "github.com/alirezadp10/hokm/internal/database"
    "gopkg.in/telebot.v4"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "log"
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

func AddPlayerToGame(db *gorm.DB, username, gameId string) (*database.Game, error) {
    var player database.Player
    db.First(&player, "username = ?", username)

    newGame := database.Game{GameId: gameId, PlayerId: player.Id, CreatedAt: time.Now()}

    err := db.Create(&newGame).Error

    return &newGame, err
}

func DoesPlayerHaveAnActiveGame(db *gorm.DB, username string) (*string, bool) {
    var result struct{ GameId string }

    db.Table("players").
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

func DoesPlayerBelongsToThisGame(db *gorm.DB, username, gameId string) bool {
    var count int64

    err := db.Table("players").
        Joins("inner join games on games.player_id = players.id").
        Where("players.username = ?", username).
        Where("games.game_id = ?", gameId).
        Count(&count).Error

    if err != nil {
        log.Fatal(err)
        return false
    }

    return count > 0
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
