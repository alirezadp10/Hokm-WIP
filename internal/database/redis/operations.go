package redis

import (
    "context"
    _ "embed"
    "encoding/json"
    "fmt"
    "github.com/redis/rueidis"
    "log"
    "strconv"
)

//go:embed matchmaking.lua
var matchmakingScript string

//go:embed place-card.lua
var placeCardScript string

var gameFields = []string{
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

func Matchmaking(ctx context.Context, client rueidis.Client, cards []string, username, gameId, lastMoveTimestamps, king, kingCards string) {
    command := client.B().Eval().Script(matchmakingScript).Numkeys(2).Key("matchmaking", "game_creation").Arg(
        username,
        gameId,
        cards[0],
        cards[1],
        cards[2],
        cards[3],
        lastMoveTimestamps,
        king,
        kingCards,
    ).Build()
    _, err := client.Do(ctx, command).ToArray()
    if err != nil {
        log.Fatalf("could not execute Lua script: %v", err)
    }
}

func GetGameInformation(ctx context.Context, client rueidis.Client, gameId string) map[string]interface{} {
    result := make(map[string]interface{})

    command := client.B().Hmget().Key("game:" + gameId).Field(gameFields...).Build()
    information, err := client.Do(ctx, command).ToArray()
    if err != nil {
        log.Fatalf("could not resolve game information: %v", err)
    }

    for key, data := range information {
        var value map[string]interface{}
        err = json.Unmarshal([]byte(data.String()), &value)
        if err != nil {
            log.Fatalf("could not resolve game information: %v", err)
        }
        result[gameFields[key]] = value["Value"]
    }

    return result
}

func SetTrump(ctx context.Context, client rueidis.Client, gameId, trump, uIndex, lastMoveTimestamp string) error {
    cmds := make(rueidis.Commands, 0, 5)
    cmds = append(cmds, client.B().Hset().Key("game:"+gameId).FieldValue().FieldValue("trump", trump).Build())
    cmds = append(cmds, client.B().Hset().Key("game:"+gameId).FieldValue().FieldValue("has_king_cards_finished", "true").Build())
    cmds = append(cmds, client.B().Hset().Key("game:"+gameId).FieldValue().FieldValue("turn", uIndex).Build())
    cmds = append(cmds, client.B().Hset().Key("game:"+gameId).FieldValue().FieldValue("last_move_timestamp", lastMoveTimestamp).Build())
    cmds = append(cmds, client.B().Publish().Channel("choosing_trump").Message(gameId+"|"+trump).Build())

    for _, resp := range client.DoMulti(ctx, cmds...) {
        if err := resp.Error(); err != nil {
            log.Fatalf("could not execute Lua script: %v", err)
            return err
        }
    }

    return nil
}

// Define a struct for the parameters to reduce argument clutter
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

func PlaceCard(
        ctx context.Context,
        client rueidis.Client,
        params PlaceCardParams,
) error {
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
    cmd := client.B().Eval().Script(placeCardScript).
        Numkeys(1).
        Key(params.GameId).
        Arg(args...).Build()

    if err := client.Do(ctx, cmd).Error(); err != nil {
        // Handle error gracefully instead of logging fatal
        return fmt.Errorf("could not execute Lua script: %w", err)
    }

    return nil
}

func Subscribe(ctx context.Context, client rueidis.Client, channel string, message func(rueidis.PubSubMessage)) error {
    err := client.Receive(ctx, client.B().Subscribe().Channel(channel).Build(), message)

    if err != nil {
        log.Printf("Error in subscribing to %v channel: %v", channel, err)
        return err
    }

    return nil
}

func Unsubscribe(ctx context.Context, client rueidis.Client, channel string) {
    unsubscribeErr := client.Do(ctx, client.B().Unsubscribe().Channel(channel).Build()).Error()
    if unsubscribeErr != nil {
        log.Println("Error while unsubscribing:", unsubscribeErr)
    }
}

func RemovePlayerList(ctx context.Context, client rueidis.Client, key, username string) {
    command := client.B().Lrem().Key(key).Count(0).Element(username).Build()
    err := client.Do(ctx, command).Error()
    if err != nil {
        log.Printf("Error in removing player list: %v", err)
    }
}
