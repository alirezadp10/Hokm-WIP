package main

import (
    "context"
    "github.com/redis/rueidis"
    "log"
)

func main() {
    ctx := context.Background()
    client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
        log.Fatal("couldn't connect to redis")
    }
    client.Do(ctx, client.B().Rpush().Key("matchmaking").Element("a").Build())
    countOfWaitingPeople, _ := client.Do(ctx, client.B().Llen().Key("matchmaking").Build()).ToInt64()
    var players []string
    if countOfWaitingPeople >= 4 {
        for i := 0; i < 4; i++ {
            player, _ := client.Do(ctx, client.B().Lpop().Key("matchmaking").Build()).ToString()
            players = append(players, player)
        }
    }
}
