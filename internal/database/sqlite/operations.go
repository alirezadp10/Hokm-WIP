package sqlite

import (
    "github.com/alirezadp10/hokm/internal/database"
    "github.com/alirezadp10/hokm/internal/request"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
    "log"
    "time"
)

func SavePlayer(db *gorm.DB, user request.User) (*database.Player, error) {
    newPlayer := database.Player{
        Id:        user.ID,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Username:  user.Username,
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

func HasGameFinished(db *gorm.DB, gameId string) bool {
    var game database.Game

    db.Table("games").Where("game_id = ?", gameId).First(&game)

    return game.FinishedAt != nil
}
