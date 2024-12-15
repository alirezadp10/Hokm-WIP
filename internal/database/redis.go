package database

import (
    "context"
    "github.com/redis/rueidis"
    "log"
)

func GetNewRedisConnection() rueidis.Client {
    client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
    if err != nil {
        log.Fatal("couldn't connect to redis")
    }
    return client
}

func Matchmaking(ctx context.Context, userId string) {
    client := GetNewRedisConnection()

    luaScript := `
        local count = redis.call('LLEN', KEYS[1])  -- Get the length of the list
        local players = {}

        redis.call('RPUSH', KEYS[1], KEYS[3])  -- Add player to a queue

        if count >= 4 then
            -- If there are at least 4 players, pop 4 players from the queue
            for i = 1, 4 do
                local player = redis.call('LPOP', KEYS[1])  -- Remove the first player from the list
                table.insert(players, player)
            end

            -- Publish the list of players to a channel 
            redis.call('PUBLISH', KEYS[2], table.concat(players, ","))
        end
        
        -- Return the list of players
        return players
	`

    command := client.B().Eval().Script(luaScript).Numkeys(2).Key("matchmaking", "waiting", userId).Build()
    players, err := client.Do(ctx, command).ToArray()
    if err != nil {
        log.Fatalf("could not execute Lua script: %v", err)
    }

    if len(players) >= 4 {
        // we need to create a game for these 4 player
    }
}
