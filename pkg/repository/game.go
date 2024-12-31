package repository

import (
    "context"
    _ "embed"
    "encoding/json"
    "github.com/alirezadp10/hokm/pkg/model"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
    "log"
)

//go:embed lua/matchmaking.lua
var matchmakingScript string

type GameRepositoryContract interface {
    HasGameFinished(gameID string) (bool, error)
    DoesPlayerBelongToGame(username, gameID string) (bool, error)
    GetGameInformation(ctx context.Context, gameID string) (map[string]interface{}, error)
    DoesPlayerHaveAnActiveGame(username string) (*string, bool)
    Matchmaking(ctx context.Context, cards []string, username, gameID, lastMoveTimestamps, king, kingCards string)
    RemovePlayerFromWaitingList(ctx context.Context, key, username string)
    GetGameInf(ctx context.Context, channel string, message func(rueidis.PubSubMessage)) error
}

var _ GameRepositoryContract = &GameRepository{}

type GameRepository struct {
    sqlite gorm.DB
    redis  rueidis.Client
}

func NewGameRepository(sqliteClient *gorm.DB, redisClient *rueidis.Client) *GameRepository {
    return &GameRepository{
        sqlite: *sqliteClient,
        redis:  *redisClient,
    }
}

func (r *GameRepository) HasGameFinished(gameID string) (bool, error) {
    var game model.Game

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

func (r *GameRepository) Matchmaking(ctx context.Context, cards []string, username, gameID, lastMoveTimestamps, king, kingCards string) {
    command := r.redis.B().Eval().Script(matchmakingScript).Numkeys(2).Key("matchmaking", "game_creation").Arg(
        username,
        gameID,
        cards[0],
        cards[1],
        cards[2],
        cards[3],
        lastMoveTimestamps,
        king,
        kingCards,
    ).Build()
    _, err := r.redis.Do(ctx, command).ToArray()
    if err != nil {
        log.Fatalf("could not execute Lua script: %v", err)
    }
}

func (r *GameRepository) RemovePlayerFromWaitingList(ctx context.Context, key, username string) {
    command := r.redis.B().Lrem().Key(key).Count(0).Element(username).Build()
    err := r.redis.Do(ctx, command).Error()
    if err != nil {
        log.Printf("Error in removing player list: %v", err)
    }
}

func (r *GameRepository) GetGameInf(ctx context.Context, channel string, message func(rueidis.PubSubMessage)) error {
    err := r.redis.Receive(ctx, r.redis.B().Subscribe().Channel("game_creation").Build(), message)

    if err != nil {
        log.Printf("Error in subscribing to %v channel: %v", channel, err)
        return err
    }

    return nil
}
