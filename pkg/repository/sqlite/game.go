package sqliteRepo

import (
    "context"
    _ "embed"
    "github.com/alirezadp10/hokm/pkg/model"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
    "log"
)

var _ repository.GameRepositoryContract = &GameRepository{}

type GameRepository struct {
    sqlite gorm.DB
}

func NewGameRepository(sqliteClient *gorm.DB) *GameRepository {
    return &GameRepository{
        sqlite: *sqliteClient,
    }
}

func (r *GameRepository) GetGameInformation(ctx context.Context, gameID string) (map[string]interface{}, error) {
    //TODO implement me
    panic("implement me")
}

func (r *GameRepository) Matchmaking(ctx context.Context, cards []string, username, gameID, lastMoveTimestamps, king, kingCards string) {
    //TODO implement me
    panic("implement me")
}

func (r *GameRepository) RemovePlayerFromWaitingList(ctx context.Context, key, username string) {
    //TODO implement me
    panic("implement me")
}

func (r *GameRepository) GetGameInf(ctx context.Context, channel string, message func(rueidis.PubSubMessage)) error {
    //TODO implement me
    panic("implement me")
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
