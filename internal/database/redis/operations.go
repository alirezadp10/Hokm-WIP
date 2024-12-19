package redis

import (
    "context"
    _ "embed"
    "encoding/json"
    "github.com/redis/rueidis"
    "log"
    "strconv"
)

//go:embed matchmaking.lua
var matchmakingScript string

func Matchmaking(ctx context.Context, client rueidis.Client, cards []string, userId, gameId, lastMoveTimestamps, king, kingCards string) {
    command := client.B().Eval().Script(matchmakingScript).Numkeys(2).Key("matchmaking", "game_creation").Arg(
        userId,
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

    fields := []string{
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
        "has_king_cards_finished",
    }

    command := client.B().Hmget().Key("game:" + gameId).Field(fields...).Build()
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
        result[fields[key]] = value["Value"]
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
            return err
        }
    }

    return nil
}

func PlaceCard(ctx context.Context, client rueidis.Client, playerIndex int, gameId, card, centerCards string) error {
    cmds := make(rueidis.Commands, 0, 2)
    cmds = append(cmds, client.B().Hset().Key("game:"+gameId).FieldValue().FieldValue("center_cards", centerCards).Build())
    cmds = append(cmds, client.B().Publish().Channel("placing_card").Message(gameId+"|"+strconv.Itoa(playerIndex)+"|"+card).Build())

    for _, resp := range client.DoMulti(ctx, cmds...) {
        if err := resp.Error(); err != nil {
            return err
        }
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
