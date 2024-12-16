package redis

import (
    "context"
    _ "embed"
    "encoding/json"
    "github.com/redis/rueidis"
    "log"
)

//go:embed matchmaking.lua
var matchmakingScript string

// Matchmaking Find an open game for a player
func Matchmaking(ctx context.Context, client rueidis.Client, userId string, gameId string) {
    command := client.B().Eval().Script(matchmakingScript).Numkeys(2).Key("matchmaking", "waiting", userId, gameId).Build()
    _, err := client.Do(ctx, command).ToArray()
    if err != nil {
        log.Fatalf("could not execute Lua script: %v", err)
    }
}

func GetGamesPlayers(ctx context.Context, client rueidis.Client, gameId string) []string {
    command := client.B().Hget().Key("game:" + gameId).Field("players").Build()
    players, err := client.Do(ctx, command).ToString()
    if err != nil {
        log.Fatalf("could not resolve players: %v", err)
    }

    var playersList []string

    _ = json.Unmarshal([]byte(players), &playersList)

    return playersList
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
