package redisRepo

import (
    "context"
    _ "embed"
    "fmt"
    "github.com/alirezadp10/hokm/pkg/repository"
    "github.com/redis/rueidis"
    "log"
    "strconv"
)

//go:embed lua/place-card.lua
var placeCardScript string

var _ repository.CardsRepositoryContract = &CardsRepository{}

type CardsRepository struct {
    redis rueidis.Client
}

func NewCardsRepository(redisClient *rueidis.Client) *CardsRepository {
    return &CardsRepository{
        redis: *redisClient,
    }
}

func (r *CardsRepository) SetTrump(ctx context.Context, gameID, trump, uIndex, lastMoveTimestamp string) error {
    cmds := make(rueidis.Commands, 0, 5)
    cmds = append(cmds, r.redis.B().Hset().Key("game:"+gameID).FieldValue().FieldValue("trump", trump).Build())
    cmds = append(cmds, r.redis.B().Hset().Key("game:"+gameID).FieldValue().FieldValue("has_king_cards_finished", "true").Build())
    cmds = append(cmds, r.redis.B().Hset().Key("game:"+gameID).FieldValue().FieldValue("turn", uIndex).Build())
    cmds = append(cmds, r.redis.B().Hset().Key("game:"+gameID).FieldValue().FieldValue("last_move_timestamp", lastMoveTimestamp).Build())
    cmds = append(cmds, r.redis.B().Publish().Channel("choosing_trump").Message(gameID+"|"+trump).Build())

    for _, resp := range r.redis.DoMulti(ctx, cmds...) {
        if err := resp.Error(); err != nil {
            log.Fatalf("could not execute Lua script: %v", err)
            return err
        }
    }

    return nil
}

func (r *CardsRepository) PlaceCard(ctx context.Context, params repository.PlaceCardParams) error {
    // Prepare the arguments for the Lua script, ensuring they are of type []string
    args := []string{
        params.CenterCards,
        params.LeadSuit,
        params.CardsWinner,
        params.Points,
        params.Turn,
        params.King,
        params.WasKingChanged,
        //TODO must be fixed
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
    cmd := r.redis.B().Eval().Script(placeCardScript).Numkeys(1).Key(params.GameId).Arg(args...).Build()

    if err := r.redis.Do(ctx, cmd).Error(); err != nil {
        // Handle error gracefully instead of logging fatal
        return fmt.Errorf("could not execute Lua script: %w", err)
    }

    return nil
}
