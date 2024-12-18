local count = redis.call('LLEN', KEYS[1])  -- Get the length of the list
local players = {}

redis.call('RPUSH', KEYS[1], ARGS[1])  -- Add player to a queue

if count >= 4 then
    -- If there are at least 4 players, pop 4 players from the queue
    for i = 1, 4 do
        local player = redis.call('LPOP', KEYS[1])  -- Remove the first player from the list
        table.insert(players, player)
    end

    -- Create a game for selected users
    redis.call('HSET', 'game:' .. ARGS[2], 'players', table.concat(players, ","))

    redis.call('HSET', 'game:' .. ARGS[2], 'center_cards', '0,0,0,0')

    -- Set players cards
    redis.call('HSET', 'game:' .. gameID, 'cards', cjson.encode({
        0 = {ARGS[3]},
        1 = {ARGS[4]}
        2 = {ARGS[5]}
        3 = {ARGS[6]}
    }))

    local current_time = os.time()

    redis.call('HSET', 'game:' .. gameID, 'last_move_timestamp', current_time)

    -- Publish the list of players to a channel
    redis.call('PUBLISH', KEYS[2], table.concat(players, ","), ARGS[2])
end

-- Return the list of players
return players