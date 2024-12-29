package repository

import (
    "context"
    _ "embed"
    "fmt"
    "github.com/redis/rueidis"
    "gorm.io/gorm"
    "log"
    "strconv"
)

//go:embed lua/place-card.lua
var placeCardScript string

type CardsRepositoryContract interface {
    SetTrump(ctx context.Context, gameID, trump, uIndex, lastMoveTimestamp string) error
    PlaceCard(ctx context.Context, params PlaceCardParams) error
}

var _ CardsRepositoryContract = &CardsRepository{}

type CardsRepository struct {
    Sqlite *gorm.DB
    Redis  rueidis.Client
}

func NewCardsRepository(sqliteClient *gorm.DB, redisClient rueidis.Client) *CardsRepository {
    return &CardsRepository{
        Sqlite: sqliteClient,
        Redis:  redisClient,
    }
}

func (r *CardsRepository) SetTrump(ctx context.Context, gameID, trump, uIndex, lastMoveTimestamp string) error {
    cmds := make(rueidis.Commands, 0, 5)
    cmds = append(cmds, r.Redis.B().Hset().Key("game:"+gameID).FieldValue().FieldValue("trump", trump).Build())
    cmds = append(cmds, r.Redis.B().Hset().Key("game:"+gameID).FieldValue().FieldValue("has_king_cards_finished", "true").Build())
    cmds = append(cmds, r.Redis.B().Hset().Key("game:"+gameID).FieldValue().FieldValue("turn", uIndex).Build())
    cmds = append(cmds, r.Redis.B().Hset().Key("game:"+gameID).FieldValue().FieldValue("last_move_timestamp", lastMoveTimestamp).Build())
    cmds = append(cmds, r.Redis.B().Publish().Channel("choosing_trump").Message(gameID+"|"+trump).Build())

    for _, resp := range r.Redis.DoMulti(ctx, cmds...) {
        if err := resp.Error(); err != nil {
            log.Fatalf("could not execute Lua script: %v", err)
            return err
        }
    }

    return nil
}

type PlaceCardParams struct {
    GameId            string
    Card              string
    CenterCards       string
    LeadSuit          string
    CardsWinner       string
    Points            string
    Turn              string
    King              string
    WasKingChanged    string
    LastMoveTimestamp string
    Trump             string
    IsItNewRound      string
    Cards             []string
    PlayerIndex       int
}

func (r *CardsRepository) PlaceCard(ctx context.Context, params PlaceCardParams) error {
    // Prepare the arguments for the Lua script, ensuring they are of type []string
    args := []string{
        params.CenterCards,
        params.LeadSuit,
        params.CardsWinner,
        params.Points,
        params.Turn,
        params.King,
        params.WasKingChanged,
        params.Cards[0],
        params.Cards[1],
        params.Cards[2],
        params.Cards[3],
        strconv.Itoa(params.PlayerIndex),
        params.Card,
        params.LastMoveTimestamp,
        params.Trump,
        params.IsItNewRound,
    }

    // Create and execute the Lua script
    cmd := r.Redis.B().Eval().Script(placeCardScript).
        Numkeys(1).
        Key(params.GameId).
        Arg(args...).Build()

    if err := r.Redis.Do(ctx, cmd).Error(); err != nil {
        // Handle error gracefully instead of logging fatal
        return fmt.Errorf("could not execute Lua script: %w", err)
    }

    return nil
}
