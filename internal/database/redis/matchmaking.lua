local count = redis.call('LLEN', KEYS[1])  -- Get the length of the list
local players = {}

redis.call('RPUSH', KEYS[1], KEYS[3])  -- Add player to a queue

if count >= 4 then
    -- If there are at least 4 players, pop 4 players from the queue
    for i = 1, 4 do
        local player = redis.call('LPOP', KEYS[1])  -- Remove the first player from the list
        table.insert(players, player)
    end

    -- Create a game for selected users
    redis.call('HSET', 'game:' .. KEYS[4], 'players', table.concat(players, ","), 'total+')

    -- Publish the list of players to a channel
    redis.call('PUBLISH', KEYS[2], table.concat(players, ","), KEYS[4])
end

-- Return the list of players
return players