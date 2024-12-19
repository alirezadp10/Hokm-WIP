redis.call('RPUSH', KEYS[1], ARGV[1])  -- Add player to a queue

local players = {}

local count = redis.call('LLEN', KEYS[1])  -- Get the length of the list

if count >= 4 then
    -- If there are at least 4 players, pop 4 players from the queue
    for i = 1, 4 do
        local player = redis.call('LPOP', KEYS[1])  -- Remove the first player from the list
        table.insert(players, player)
    end

    -- Create a game for selected users
    redis.call('HSET', 'game:' .. ARGV[2], 'players', table.concat(players, ","))

    redis.call('HSET', 'game:' .. ARGV[2], 'points', '{"total":"0,0","round":"0,0"}')

    redis.call('HSET', 'game:' .. ARGV[2], 'center_cards', ',,,')

    -- Set players cards
    redis.call('HSET', 'game:' .. ARGV[2], 'cards', cjson.encode({
        [0] = cjson.decode(ARGV[3]),
        [1] = cjson.decode(ARGV[4]),
        [2] = cjson.decode(ARGV[5]),
        [3] = cjson.decode(ARGV[6])
    }))

    redis.call('HSET', 'game:' .. ARGV[2], 'trump', '')

    redis.call('HSET', 'game:' .. ARGV[2], 'turn', '')

    redis.call('HSET', 'game:' .. ARGV[2], 'last_move_timestamp', ARGV[7])

    redis.call('HSET', 'game:' .. ARGV[2], 'judge', ARGV[8])

    redis.call('HSET', 'game:' .. ARGV[2], 'judge_cards', ARGV[9])

    redis.call('HSET', 'game:' .. ARGV[2], 'has_judge_cards_finished', "false")

    -- Publish the list of players to a channel
    redis.call('PUBLISH', KEYS[2], table.concat(players, ",") .. "|" .. ARGV[2])
end

-- Return the list of players
return players