package redis

import (
    "context"
    _ "embed"
    "encoding/json"
    "github.com/redis/rueidis"
    "log"
    "strings"
)

//go:embed matchmaking.lua
var matchmakingScript string

// Matchmaking Find an open game for a player
func Matchmaking(ctx context.Context, client rueidis.Client, userId string, gameId string, cards [][]string) {
    command := client.B().Eval().
        Script(matchmakingScript).
        Numkeys(2).
        Key("matchmaking", "game_creation").
        Arg(userId, gameId, strings.Join(cards[0], ","), strings.Join(cards[1], ","), strings.Join(cards[2], ","), strings.Join(cards[3], ",")).
        Build()
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
        "judge",
        "trump",
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

func SetTrump(ctx context.Context, client rueidis.Client, gameId, trump string) error {
    err := client.Do(ctx, client.B().Hset().Key("game:"+gameId).FieldValue().FieldValue("trump", trump).Build()).Error()
    if err != nil {
        return err
    }
    return nil
}
