package main

import (
    "context"
    "fmt"
    "github.com/alirezadp10/hokm/internal/database/redis"
    "github.com/alirezadp10/hokm/internal/hokm"
    "github.com/alirezadp10/hokm/internal/utils/my_slice"
    "github.com/google/uuid"
    "github.com/joho/godotenv"
    "github.com/redis/rueidis"
    "log"
    "strconv"
    "strings"
)

func main() {
    _ = godotenv.Load()

    client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
        log.Fatalf("could not connect to Redis: %v", err)
    }
    defer client.Close()

    ctx := context.Background()
    gameId := uuid.New().String()
    distributedCards := hokm.DistributeCards()
    redis.Matchmaking(ctx, client, "1", gameId, distributedCards)
    redis.Matchmaking(ctx, client, "2", gameId, distributedCards)
    redis.Matchmaking(ctx, client, "3", gameId, distributedCards)
    redis.Matchmaking(ctx, client, "4", gameId, distributedCards)

    gameInformation := redis.GetGameInformation(ctx, client, gameId)
    players := strings.Split(gameInformation["players"].(string), ",")
    uIndex := my_slice.GetIndex("1", players)

    judgeIndex, _ := strconv.Atoi(gameInformation["judge"].(string))
    turnIndex, _ := strconv.Atoi(gameInformation["turn"].(string))

    fmt.Println(map[string]interface{}{
        "players":      hokm.GetPlayersWithDirections(players, uIndex),
        "points":       hokm.GetPoints(gameInformation["points"].(string), uIndex),
        "centerCards":  hokm.GetCenterCards(gameInformation["center_cards"].(string), uIndex),
        "turn":         hokm.GetDirection(turnIndex, uIndex),
        "judge":        hokm.GetDirection(judgeIndex, uIndex),
        "timeRemained": hokm.GetTimeRemained(gameInformation["last_move_timestamp"].(string)),
        "kingsCards":   hokm.GetKingsCards(gameInformation["kings_cards"].(string), uIndex),
        "yourCards":    hokm.GetYourCards(gameInformation["cards"].(string), uIndex),
        "trump":        gameInformation["trump"],
    })
}
