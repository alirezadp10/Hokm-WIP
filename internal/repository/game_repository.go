package repository

import (
    "context"
    "encoding/json"
    "github.com/alirezadp10/hokm/internal/database"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
    "log"
)

type GameRepository struct {
    sqlite *gorm.DB
    redis  rueidis.Client
}

func (r *GameRepository) HasGameFinished(gameID string) (bool, error) {
    var game database.Game

    r.sqlite.Table("games").Where("game_id = ?", gameID).First(&game)

    return game.FinishedAt != nil, nil
}

func (r *GameRepository) DoesPlayerBelongToGame(username, gameID string) (bool, error) {
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

func (r *GameRepository) GetGameInformation(ctx context.Context, gameID string) (map[string]interface{}, error) {
    result := make(map[string]interface{})

    gameFields := []string{
        "players",
        "points",
        "center_cards",
        "current_turn",
        "players_cards",
        "king",
        "trump",
        "turn",
        "last_move_timestamp",
        "cards",
        "king_cards",
        "lead_suit",
        "has_king_cards_finished",
        "who_has_won_the_cards",
        "who_has_won_the_round",
        "who_has_won_the_game",
        "was_king_changed",
    }

    command := r.redis.B().Hmget().Key("game:" + gameID).Field(gameFields...).Build()
    information, err := r.redis.Do(ctx, command).ToArray()
    if err != nil {
        log.Fatalf("could not resolve game information: %v", err)
        return nil, err
    }

    for key, data := range information {
        var value map[string]interface{}
        err = json.Unmarshal([]byte(data.String()), &value)
        if err != nil {
            log.Fatalf("could not resolve game information: %v", err)
            return nil, err
        }
        result[gameFields[key]] = value["Value"]
    }

    return result, nil
}

func (r *GameRepository) DoesPlayerHaveAnActiveGame(username string) (*string, bool) {
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
